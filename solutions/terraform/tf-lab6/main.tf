terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

module "static_website_bucket" {
  source = "./modules/gcp-storage-static-website-bucket"

  project_id  = var.project_id
  bucket_name = "jrs-example-2023-01-15"
  location    = "US"
  admin_email = "alexander.nolan@sudocloud.io"

  labels = {
    environment = "dev"
    terraform   = "true"
  }
}