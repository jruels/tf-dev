terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }
}

provider "google" {
  project = "iis-tf-dev"
  region  = "us-central1"
}

resource "google_compute_instance" "tf-example-import" {
  count        = 3
  name         = "tf-example-import-${count.index}"
  machine_type = "e2-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral public IP
    }
  }

  labels = {
    role = "terraform"
    name = "tf-example-import-${count.index}"
  }
} 