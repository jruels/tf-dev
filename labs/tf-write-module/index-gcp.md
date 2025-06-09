# Write a Terraform module

## Overview
In this lab, you will create a module to manage Google Cloud Storage buckets to host static websites.

## Module structure

Remember the typical structure for a new module is: 
```
.
├── LICENSE
├── README.md
├── main.tf
├── variables.tf
├── outputs.tf
```
None of these files are required, or have any special meaning to Terraform when it uses your module. You can create a module with a single `.tf` file, or use any other file structure you like.

Each of these files serves a purpose:

- `LICENSE` will contain the license under which your module will be distributed. When you share your module, the `LICENSE` file will let people using it know the terms under which it has been made available. Terraform itself does not use this file.
- `README.md` will contain documentation describing how to use your module, in markdown format. Terraform does not use this file, but services like the Terraform Registry and GitHub will display the contents of this file to people who visit your module's Terraform Registry or GitHub page.
- `main.tf` will contain the main set of configuration for your module. You can also create other configuration files and organize them however makes sense for your project.
- `variables.tf` will contain the variable definitions for your module. When your module is used by others, the variables will be configured as arguments in the `module block`. Since all Terraform values must be defined, any variables that are not given a default value will become required arguments. Variables with default values can also be provided as module arguments, overriding the default value.
- `outputs.tf` will contain the output definitions for your module. Module outputs are made available to the configuration using the module, so they are often used to pass information about the parts of your infrastructure defined by the module to other parts of your configuration.

You also want to make sure and add the following to your ignore list. If you are tracking your module in GitHub use `.gitignore`

- `terraform.tfstate` and `terraform.tfstate.backup`: These files contain your Terraform state, and are how Terraform keeps track of the relationship between your configuration and the infrastructure provisioned by it.
- `.terraform`: This directory contains the modules and plugins used to provision your infrastructure. These files are specific to a specific instance of Terraform when provisioning infrastructure, not the configuration of the infrastructure defined in `.tf` files.
- `*.tfvars`: Since module input variables are set via arguments to the module block in your configuration, you don't need to distribute any `*.tfvars` files with your module, unless you are also using it as a standalone Terraform configuration.

## Create a module 

### Create Terraform configuration

1. In **Visual Studio Code**, open the working directory created in the previous lab (`YYYYMMDD/terraform`).
2. Right-click in the **Explorer** pane and select **New Folder**.
3. Name the folder `tf-lab6`.
4. Create a sub-folder called `modules` inside `tf-lab6`.
5. Inside `modules`, create a sub-folder named `gcp-storage-static-website-bucket`.

After creating these directories, your configuration's directory structure will look like this:
```
tf-lab6
├── LICENSE
├── README.md
├── main.tf
├── modules
│   └── gcp-storage-static-website-bucket
├── outputs.tf
└── variables.tf
```

Hosting a static website with Cloud Storage is a fairly common use case. While it isn't too difficult to figure out the correct configuration to provision a bucket this way, encapsulating this configuration within a module will provide your users with a quick and easy way to create buckets they can use to host static websites that adhere to best practices. Another benefit of using a module is that the module name can describe exactly what buckets created with it are for. In this example, the `gcp-storage-static-website-bucket` module creates Cloud Storage buckets that host static websites.

You will work with three Terraform configuration files inside the `gcp-storage-static-website-bucket` directory: `main.tf`, `variables.tf`, and `outputs.tf`.

Inside the `modules/gcp-storage-static-website-bucket` directory, create a `main.tf` with the following: 

```hcl
resource "google_storage_bucket" "website" {
  name          = var.bucket_name
  location      = var.location
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD", "OPTIONS"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  labels = var.labels
}

resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.website.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_storage_bucket_iam_member" "admin" {
  bucket = google_storage_bucket.website.name
  role   = "roles/storage.objectAdmin"
  member = "user:${var.admin_email}"
}
```

This configuration creates a public Cloud Storage bucket configured for website hosting with an index page and a 404 error page.

You will notice that there is no provider block in this configuration. When Terraform processes a module block, it will inherit the provider from the enclosing configuration. Because of this, there's no need to include `provider` blocks in modules.

Define the following variables in `variables.tf` inside the `modules/gcp-storage-static-website-bucket` directory:

```hcl
variable "project_id" {
  description = "The ID of the project where the bucket will be created"
  type        = string
}

variable "bucket_name" {
  description = "Name of storage bucket. Must be unique"
  type        = string
}

variable "location" {
  description = "Location for the storage bucket"
  type        = string
  default     = "US"
}

variable "admin_email" {
  description = "Email address of the bucket admin"
  type        = string
}

variable "labels" {
  description = "Labels to set on the bucket"
  type        = map(string)
  default     = {}
}
```

When using a module, variables are set by passing arguments to the module in your configuration. You will set values for these variables when calling this module from your root module's `main.tf`.

Inside the `modules/gcp-storage-static-website-bucket` directory, add outputs to your module in the `outputs.tf` file:

```hcl
output "url" {
  description = "The URL of the website"
  value       = "https://storage.googleapis.com/${google_storage_bucket.website.name}/index.html"
}

output "name" {
  description = "Name of the bucket"
  value       = google_storage_bucket.website.name
}

output "self_link" {
  description = "Self link of the bucket"
  value       = google_storage_bucket.website.self_link
}
```

Like variables, outputs in modules perform the same function as they do in the root module but are accessed differently. A module's outputs can be accessed as read-only attributes on the module object, which is available within the configuration that calls the module. You can reference these outputs in expressions as `module.<MODULE NAME>.<OUTPUT NAME>`.

Now that you have created your module, create a `main.tf` in your root module directory and add a reference to the new module:

```hcl
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
  bucket_name = "<UNIQUE BUCKET NAME>"
  location    = "US"
  admin_email = "your.email@example.com"

  labels = {
    environment = "dev"
    terraform   = "true"
  }
}
```

Google Cloud Storage buckets must be globally unique. Because of this, you will need to **replace `<UNIQUE BUCKET NAME>`** with a unique, valid name for a storage bucket. Using your name and the date is usually a good way to create a unique bucket name. For example:

```hcl
  bucket_name = "jrs-example-2023-01-15"
```

Create a `variables.tf` file in your root module directory:

```hcl
variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "Google Cloud region"
  type        = string
  default     = "us-central1"
}
```

## Define outputs
Earlier, you added several outputs to the `gcp-storage-static-website-bucket` module, making those values available to your root module configuration.

Add these values as outputs to your root module by adding the following to `outputs.tf` file in your root module directory:

```hcl
output "website_url" {
  description = "URL of the website"
  value       = module.static_website_bucket.url
}

output "bucket_name" {
  description = "Name of the bucket"
  value       = module.static_website_bucket.name
}

output "bucket_self_link" {
  description = "Self link of the bucket"
  value       = module.static_website_bucket.self_link
}
```

## Install the local module
When you add a new module to a configuration, Terraform must install it before it can be used. Both the `terraform get` and `terraform init` commands will install and update modules. The `terraform init` command will also initialize backends and install plugins.

```sh
terraform init
```

Create a `terraform.tfvars` file with your project ID:

```hcl
project_id = "YOUR_PROJECT_ID"
```

## Apply the configuration

Now you can apply the configuration:

```sh
terraform apply
```

After the configuration is applied, you'll see the outputs including the website URL, bucket name, and bucket self link.

## Cleanup

When you're done, clean up the resources:

```sh
terraform destroy
```

## Congratulations

You have successfully created a Terraform module for hosting static websites using Google Cloud Storage! 