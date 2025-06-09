terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "random_pet" "petname" {
  length    = 3
  separator = "-"
}

resource "google_storage_bucket" "website" {
  name          = "${var.environment}-${random_pet.petname.id}"
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  labels = {
    environment = var.environment
  }
}

resource "google_storage_bucket_iam_member" "website" {
  bucket = google_storage_bucket.website.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_storage_bucket_object" "website" {
  name          = "index.html"
  bucket        = google_storage_bucket.website.name
  content       = file("${path.module}/assets/index.html")
  content_type  = "text/html"
}