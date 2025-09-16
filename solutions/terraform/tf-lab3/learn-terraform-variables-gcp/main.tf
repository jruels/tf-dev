terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.1.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "~> 7.0"

  project_id   = var.project_id
  network_name = "${var.project_name}-${var.environment}"
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name   = "private-subnet-0"
      subnet_ip     = var.private_subnet_cidr_blocks[0]
      subnet_region = var.region
      secondary_ip_ranges = [
        {
          range_name    = "secondary-range-0"
          ip_cidr_range = slice(var.private_subnet_0_secondary_ranges, 0, 3)[0]
        },
        {
          range_name    = "secondary-range-1"
          ip_cidr_range = slice(var.private_subnet_0_secondary_ranges, 0, 3)[1]
        },
        {
          range_name    = "secondary-range-2"
          ip_cidr_range = slice(var.private_subnet_0_secondary_ranges, 0, 3)[2]
        }
      ]
    },
    {
      subnet_name   = "private-subnet-1"
      subnet_ip     = var.private_subnet_cidr_blocks[1]
      subnet_region = var.region
      secondary_ip_ranges = [
        {
          range_name    = "subnet-1-secondary-0"
          ip_cidr_range = var.private_subnet_1_secondary_ranges[0]
        },
        {
          range_name    = "subnet-1-secondary-1"
          ip_cidr_range = var.private_subnet_1_secondary_ranges[1]
        }
      ]
    },
    {
      subnet_name   = "public-subnet-0"
      subnet_ip     = var.public_subnet_cidr_blocks[0]
      subnet_region = var.region
    },
    {
      subnet_name   = "public-subnet-1"
      subnet_ip     = var.public_subnet_cidr_blocks[1]
      subnet_region = var.region
    }
  ]

  routes = [
    {
      name              = "egress-internet"
      description       = "Route through IGW to access internet"
      destination_range = "0.0.0.0/0"
      next_hop_internet = "true"
    }
  ]
}


module "cloud_router" {
  source  = "terraform-google-modules/cloud-router/google"
  version = "~> 5.0"

  name    = "nat-router"
  project = var.project_id
  region  = var.region
  network = module.vpc.network_name

  nats = [{
    name                               = "nat-config"
    source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
    subnetworks = [
      {
        name                     = module.vpc.subnets["${var.region}/private-subnet-0"].self_link
        source_ip_ranges_to_nat  = ["ALL_IP_RANGES"]
        secondary_ip_range_names = []
      },
      {
        name                     = module.vpc.subnets["${var.region}/private-subnet-1"].self_link
        source_ip_ranges_to_nat  = ["ALL_IP_RANGES"]
        secondary_ip_range_names = []
      }
    ]
  }]
}

# Create individual instances
module "instances" {
  source = "./modules/gcp-instance"

  instance_count = var.instance_count
  name_prefix    = "${var.project_name}-${var.environment}"
  machine_type   = var.machine_type
  zone          = var.zone
  network       = module.vpc.network_name
  subnetwork    = module.vpc.subnets["${var.region}/private-subnet-0"].name
  network_tags   = var.network_tags
  labels        = var.resource_labels
}

# Create an unmanaged instance group for the load balancer
resource "google_compute_instance_group" "webservers" {
  name        = "webserver-group"
  description = "Web server instance group"
  zone        = var.zone
  network     = module.vpc.network_self_link

  instances = module.instances.instance_self_links

  named_port {
    name = "http"
    port = 80
  }

}

# Update the load balancer to use the unmanaged instance group
module "lb" {
  source  = "GoogleCloudPlatform/lb-http/google"
  version = "~> 9.0"

  project           = var.project_id
  name              = "lb-${random_string.lb_id.result}"
  target_tags       = var.network_tags
  firewall_networks = [module.vpc.network_name]
  create_ipv6_address = var.create_ipv6_address

  labels = var.resource_labels

  backends = {
    default = {
      description                     = null
      protocol                       = "HTTP"
      port                           = 80
      port_name                      = "http"
      timeout_sec                    = 10
      connection_draining_timeout_sec = null
      enable_cdn                     = false
      security_policy                = null
      custom_request_headers         = null
      custom_response_headers        = null
      compression_mode               = null

      health_check = {
        check_interval_sec  = 10
        timeout_sec         = 5
        healthy_threshold   = 2
        unhealthy_threshold = 3
        request_path        = "/index.html"
        port               = 80
        host               = null
        logging            = null
      }

      log_config = {
        enable = false
        sample_rate = null
      }

      groups = [
        {
          group = google_compute_instance_group.webservers.self_link
        }
      ]

      iap_config = {
        enable               = false
        oauth2_client_id     = null
        oauth2_client_secret = null
      }
    }
  }
}

# Random string for unique names
resource "random_string" "lb_id" {
  length  = 3
  special = false
  lower   = true
  upper   = false
  numeric = true
}
