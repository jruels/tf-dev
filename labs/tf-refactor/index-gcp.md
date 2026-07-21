# Terraform Refactor codebase

## Overview 
In this lab, you will refactor a monolithic Terraform configuration into a more maintainable design. 

You will start with two copies of a web application hosted in Cloud Storage buckets (one for production, one for development). Then you will refactor this to use a single set of resources with environment-specific variables.

## Start with a monolith configuration
### Create the Lab Directory

1. In **Visual Studio Code**, open the working directory created in the previous lab (`YYYYMMDD/terraform`).
2. Right-click in the **Explorer** pane and select **New Folder**.
3. Name the folder `tf-lab7`.
4. Create the following files in your `tf-lab7` directory:

`main.tf` - configures the resources that make up your infrastructure:

```hcl
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

# Production bucket
resource "google_storage_bucket" "prod" {
  name          = "${var.prod_prefix}-${random_pet.petname.id}"
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}

resource "google_storage_bucket_object" "prod" {
  name          = "index.html"
  bucket        = google_storage_bucket.prod.name
  content       = file("${path.module}/assets/index.html")
  content_type  = "text/html"
}

# Development bucket
resource "google_storage_bucket" "dev" {
  name          = "${var.dev_prefix}-${random_pet.petname.id}"
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}

resource "google_storage_bucket_object" "dev" {
  name          = "index.html"
  bucket        = google_storage_bucket.dev.name
  content       = file("${path.module}/assets/index.html")
  content_type  = "text/html"
}
```

`variables.tf` - declares input variables:

```hcl
variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "Region for GCP resources"
  type        = string
  default     = "us-central1"
}

variable "prod_prefix" {
  description = "Prefix for production resources"
  type        = string
  default     = "prod"
}

variable "dev_prefix" {
  description = "Prefix for development resources"
  type        = string
  default     = "dev"
}
```

`terraform.tfvars` - defines your variables:

```hcl
project_id  = "YOUR_PROJECT_ID"
region      = "us-central1"
prod_prefix = "prod"
dev_prefix  = "dev"
```

`outputs.tf` - specifies the website endpoints:

```hcl
output "prod_website_url" {
  description = "URL of the production website"
  value       = "https://storage.googleapis.com/${google_storage_bucket.prod.name}/index.html"
}

output "dev_website_url" {
  description = "URL of the development website"
  value       = "https://storage.googleapis.com/${google_storage_bucket.dev.name}/index.html"
}
```

Create an `assets` directory and add an `index.html` file with some sample content:

```html
<!DOCTYPE html>
<html>
<head>
    <title>My Static Website</title>
</head>
<body>
    <h1>Welcome to my website!</h1>
</body>
</html>
```

Initialize and apply your Terraform configuration:

```sh
terraform init
terraform apply
```

Navigate to the web addresses from the Terraform output to display the deployments in a browser.

Clean up the monolithic resources before refactoring:
```bash
terraform destroy -auto-approve
```

## Refactor the configuration

Having duplicate resources for each environment creates maintenance overhead and potential inconsistencies. Let's refactor to use a single set of resources with environment-specific variables.

1. Update `variables.tf` to use environment variables:

```hcl
variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "Region for GCP resources"
  type        = string
  default     = "us-central1"
}

variable "environment" {
  description = "Environment (dev or prod)"
  type        = string
}
```

2. Create environment-specific variable files:

`dev.tfvars`:
```hcl
project_id         = "YOUR_PROJECT_ID"
region            = "us-central1"
environment       = "dev"
```

`prod.tfvars`:
```hcl
project_id         = "YOUR_PROJECT_ID"
region            = "us-central1"
environment       = "prod"
```

3. Update `main.tf` to use a single set of resources:

```hcl
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
```

4. Update `outputs.tf` to use the single resource:

```hcl
output "website_url" {
  description = "URL of the website"
  value       = "https://storage.googleapis.com/${google_storage_bucket.website.name}/index.html"
}
```

## Deploy the environments

First, let's look at what would be created in each environment:

```bash
terraform plan -var-file="dev.tfvars"
```

Notice how the bucket will be created with the `environment = "dev"` label.

Now check the production plan:
```bash
terraform plan -var-file="prod.tfvars"
```

Notice the difference in the bucket's label where `environment = "prod"`.



## Important Note About State Management

While using different tfvars files makes the code more maintainable, it introduces challenges when running Terraform locally:

1. You need to remember to destroy one environment before creating another (since they share the same state file)
2. There's risk of accidentally applying the wrong tfvars file
3. Multiple team members could try to manage different environments simultaneously

This is why in real-world scenarios, you should:
1. Run Terraform through CI/CD pipelines
2. Use separate state files for each environment
3. Implement proper state locking

These topics will be covered in the next lab on CI/CD integration.

## Cleanup

Clean up the resources:
```bash
terraform destroy -var-file="prod.tfvars"
```

## Congratulations!

You have successfully:
1. Started with a monolithic configuration that had duplicate resources
2. Refactored it to use a single set of resources with environment-specific variables
3. Learned how to deploy to different environments using tfvars files

## Problems with this Approach

While the above approach works, it has several drawbacks:

1. **Code Duplication**: Having separate directories means duplicating code between environments, which can lead to:
   - Maintenance overhead
   - Inconsistencies between environments
   - Higher chance of errors when making changes

2. **Harder to Add Environments**: Adding a new environment (like staging) requires:
   - Creating a new directory
   - Copying all files
   - Maintaining another set of identical code

3. **No Single Source of Truth**: When infrastructure code exists in multiple places, it's harder to:
   - Ensure all environments are using the same resource configurations
   - Track changes across environments
   - Maintain consistency

## A Better Approach

A better way to manage multiple environments is to:

1. Keep a single set of Terraform configurations
2. Use environment-specific `.tfvars` files
3. Use workspace or backend configurations for state management

Example:

Instead of separate directories, your structure should look like:
```
.
├── assets/
│   └── index.html
├── main.tf           # Single source of truth for all environments
├── variables.tf      # All variable definitions
├── outputs.tf        # All outputs
├── dev.tfvars       # Development environment values
└── prod.tfvars      # Production environment values
```

Then deploy to different environments using:
```bash
terraform apply -var-file="dev.tfvars"   # For development
terraform apply -var-file="prod.tfvars"  # For production
```

Benefits:
1. Single source of truth for infrastructure code
2. Easy to add new environments (just add a new .tfvars file)
3. Guaranteed consistency between environments
4. Easier to maintain and update
5. Better version control and change tracking

## Next Steps

In a real-world scenario, you would:
1. Use different state files for each environment (using backend configuration)
2. Implement proper state locking
3. Use CI/CD pipelines for deployments
4. Add environment-specific access controls 