# Terraform - use modules

## Overview 
In this lab, you will use modules from the Terraform Registry to provision a Google Cloud environment. The concepts you use in this tutorial will apply to any modules from any source.

## The Terraform Registry
Open the [Terraform Registry page for the GCP Network module](https://registry.terraform.io/modules/terraform-google-modules/network/google/latest) in your browser.

You will see information about the module and a link to the source repository. On the right side of the page, you will see a drop-down interface to select the module version and instructions for using it to provision infrastructure.

When calling a module, the `source` argument is required. In this example, Terraform will search for a module that matches the given string in the Terraform registry. You could also use a URL or local file path for the source of your modules.

The other argument shown here is the `version`. For supported sources, the version will let you define which version or versions of the module will be loaded. In this lab, you will specify the exact version number of the modules you use. 

Other arguments to module blocks are treated as input variables to the modules.

## Create Terraform configuration
Now use modules to create an example GCP environment using a Virtual Private Cloud (VPC) network and two Compute Engine instances.

1. In **Visual Studio Code**, open the working directory, `terraform`.
2. Right-click in the **Explorer** pane and select **New Folder**.
3. Name the folder `tf-lab5`.

Create a new file called `main.tf` with the following content:

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

module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "~> 7.0"

  project_id   = var.project_id
  network_name = var.vpc_name
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name   = "subnet-01"
      subnet_ip     = var.subnet_primary_cidr
      subnet_region = var.region
    }
  ]

  secondary_ranges = {
    subnet-01 = []
  }
}
```

## Set values for module input variables
You must pass input variables to the module configuration to use most modules. The configuration that calls a module is responsible for setting its input values, which are passed as arguments in the module block. Aside from `source` and `version`, most of the arguments to a module block will set variable values.

On the Terraform registry page for the GCP Network module, you will see an `Inputs` tab that describes all the input variables that the module supports.

Some input variables are required, meaning that the module doesn't provide a default value—an explicit value must be provided for Terraform to run correctly.

Review the input variables you are setting within the module `vpc` block:
- `project_id` is your Google Cloud project ID
- `network_name` will be the name of the VPC network
- `routing_mode` defines the network-wide routing mode
- `subnets` defines the subnets within your VPC network
- `secondary_ranges` defines secondary IP ranges for your subnets (empty in this example)

## Define root input variables
Create the following in `variables.tf`:

```hcl
variable "project_id" {
  description = "The ID of the project where resources will be created"
  type        = string
}

variable "region" {
  description = "The region where resources will be created"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "The zone where resources will be created"
  type        = string
  default     = "us-central1-a"
}

variable "vpc_name" {
  description = "Name of VPC"
  type        = string
  default     = "example-vpc"
}

variable "subnet_primary_cidr" {
  description = "The primary CIDR range for the subnet"
  type        = string
  default     = "10.0.0.0/24"
}
```

## Define root output values
Create the following in `outputs.tf`:

```hcl
output "network_name" {
  description = "The name of the VPC being created"
  value       = module.vpc.network_name
}

output "subnets_names" {
  description = "The names of the subnets being created"
  value       = module.vpc.subnets_names
}
```

## Provision infrastructure 
Initialize the Terraform configuration to download the provider and modules:

```bash
terraform init
```

Create a `terraform.tfvars` file with your project ID:

```hcl
project_id = "YOUR_PROJECT_ID"
```

Now run `terraform apply` to create the VPC network and instances:

```bash
terraform apply
```

You will notice that many more resources than just the VPC and compute instances will be created. The modules we used define what those resources are.

You should now see the network name, subnet names, and instance IPs output to the terminal.

## Understand how modules work
When using a new module for the first time, you must run either `terraform init` or `terraform get` to install the module. When either of these commands are run, Terraform will install any new modules in the `.terraform/modules` directory within your configuration's working directory. For local modules, Terraform will create a symlink to the module's directory. Because of this, any changes to local modules will be effective immediately, without having to re-run `terraform get`.

## Cleanup

Run `terraform destroy` to remove resources:

```bash
terraform destroy
```

Remove the `.terraform` directory to free up disk space:

```bash
rm -rf .terraform
```

## Congratulations

You have successfully used Terraform modules to create a VPC network and compute instances in Google Cloud Platform! 