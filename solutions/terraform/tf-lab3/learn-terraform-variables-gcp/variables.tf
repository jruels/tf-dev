# Variable declarations

variable "project_id" {
  description = "The ID of the GCP project"
  type        = string
  default     = "iis-tf-dev"
}

variable "region" {
  description = "The region to deploy resources to"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "The zone to deploy resources to"
  type        = string
  default     = "us-central1-a"
}

variable "private_subnet_cidr_blocks" {
  description = "Available CIDR blocks for private subnets"
  type        = list(string)
  default     = [
    "10.0.101.0/24",
    "10.0.102.0/24"
  ]
}

variable "public_subnet_cidr_blocks" {
  description = "Available CIDR blocks for public subnets"
  type        = list(string)
  default     = [
    "10.0.1.0/24",
    "10.0.2.0/24"
  ]
}

variable "instance_count" {
  description = "Number of instances to provision"
  type        = number
  default     = 2
}

variable "environment" {
  description = "Environment (dev/staging/prod)"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "project-alpha"
}

variable "machine_type" {
  description = "GCP machine type"
  type        = string
  default     = "e2-micro"
}

variable "network_tags" {
  description = "Network tags to apply to instances"
  type        = list(string)
  default     = ["web-server", "allow-health-check"]
}

variable "create_ipv6_address" {
  description = "Allocate a new IPv6 address for the load balancer. Conflicts with manually specified ipv6_address."
  type        = bool
  default     = false
}

variable "secondary_ip_ranges" {
  description = "Available CIDR blocks for secondary IP ranges."
  type        = list(string)
  default     = [
    "10.0.101.0/24",
    "10.0.102.0/24",
    "10.0.103.0/24",
    "10.0.104.0/24",
    "10.0.105.0/24",
    "10.0.106.0/24",
    "10.0.107.0/24",
    "10.0.108.0/24",
  ]
}

variable "secondary_ip_range_cidrs" {
  description = "CIDR blocks for secondary IP ranges"
  type        = list(string)
  default     = [
    "192.168.10.0/24",
    "192.168.20.0/24",
    "192.168.30.0/24",
    "192.168.40.0/24",
    "192.168.50.0/24",
    "192.168.60.0/24",
    "192.168.70.0/24",
    "192.168.80.0/24"
  ]
}

variable "private_subnet_0_secondary_ranges" {
  description = "Secondary IP CIDR ranges for private subnet 0"
  type        = list(string)
  default     = [
    "192.168.10.0/24",
    "192.168.20.0/24",
    "192.168.30.0/24"
  ]
}

variable "private_subnet_1_secondary_ranges" {
  description = "Secondary IP CIDR ranges for private subnet 1"
  type        = list(string)
  default     = [
    "192.168.40.0/24",
    "192.168.50.0/24",
    "192.168.60.0/24"
  ]
}

variable "resource_labels" {
  description = "Labels to set for all resources"
  type        = map(string)
  default     = {
    project     = "project-alpha"
    environment = "dev"
  }
}