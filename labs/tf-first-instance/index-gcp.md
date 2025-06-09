# Create your first Terraform resource

## Overview

This lab walks through setting up Terraform and creating your first resources.

## Authentication
Before proceeding, authenticate with Google Cloud:

```sh
gcloud auth application-default login
```

## Setup

These instructions assume that you are using Visual Studio Code.

### Create a Working Directory

1. Open **Visual Studio Code**.
2. Click on **File** > **Open Folder...** and create a new folder in your Documents folder named with today's date (e.g., 20240329), and enter it.
3. Inside the opened folder (e.g. 20240329), click the **New Folder** icon.
4. Name the new folder `terraform` and double-click it to open it in Visual Studio Code.

## Create Terraform Configuration

### Create the Lab Directory

1. In Visual Studio Code, click the new folder icon or Right-click inside the **Explorer** pane and select **New Folder**.
2. Name this folder `tf-lab1` and expand it.

### Create the `main.tf` File

1. Right-click the `tf-lab1` folder, select **New File**, and name it `main.tf`.
2. Open `main.tf` and paste the following Terraform configuration:

```hcl
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
    }
  }
}

provider "google" {
  project = "YOUR_PROJECT_ID"
  region  = "us-west1"
}

resource "google_compute_instance" "lab1-tf-example" {
  name         = "lab1-tf-example"
  machine_type = "e2-micro"
  zone         = "us-west1-a"

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
    name = "lab1-tf-example"
  }
}
```

Replace 'YOUR_PROJECT_ID' with your actual project ID, once done this configuration is ready to be applied. We'll review each section in detail below.

## Providers

Each Terraform module must declare required providers so that Terraform can install and use them. The `provider` block configures Terraform to use the Google provider.

### Security Note

**Never hard-code credentials** into `*.tf` files. Instead, use service account credentials configured through environment variables or application default credentials.

## Initialize the Directory

1. In **Visual Studio Code**, right-click on the `tf-lab1` folder in the **Explorer** pane.
2. Select **Open in Integrated Terminal** to ensure that all Terraform commands run in the correct directory.
3. Run the following command to initialize Terraform and download the required providers:
   ```sh
   terraform init
   ```
4. You should see Terraform downloading the necessary plugins and initializing the environment.

## Format and Validate Configuration

1. Run the following in the terminal 
   ```hcl
   terraform fmt
   ```

   * This command formats the file to meet HCL standards

   ```hcl
   terraform validate
   ```

   * This command runs syntax checks and confirms all variables are declared.

## Create Infrastructure

1. Open **Command Palette** and type **Terraform: Plan**, then select it.
   - Review the planned changes before applying them.
2. Apply the configuration:
   - In the terminal, type:
     ```sh
     terraform apply
     ```
   - Terraform will prompt for confirmation. Type **yes** when prompted.

Terraform will create a Google Compute Engine instance as specified in your configuration.

## Cleanup

To delete the resources created:

1. Destroy the instance:
   - In the terminal, type:
     ```
     terraform destroy
     ```
   - Terraform will prompt for confirmation. Type **yes** when prompted.

## Congratulations!

You've successfully created and destroyed your first Terraform-managed infrastructure. 