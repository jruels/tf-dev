# Terraform CI/CD Lab 1: Environment Configuration

## Overview
In this lab, you will prepare your Terraform configuration for CI/CD by making it environment-agnostic and creating environment-specific variable files.

## Starting Point
We'll use the configuration from `tf-lab3` which deploys a VPC network with subnets, instances, and a load balancer.

## Steps

### 1. Create Environment Variable Files

1. Create `dev.tfvars`:
```hcl
project_id   = "iis-tf-dev"
environment  = "dev"
project_name = "tf-lab3"
region       = "us-central1"
zone         = "us-central1-a"

# Network configuration
private_subnet_cidr_blocks = ["10.0.1.0/24", "10.0.2.0/24"]
public_subnet_cidr_blocks  = ["10.0.3.0/24", "10.0.4.0/24"]
private_subnet_0_secondary_ranges = [
  "10.1.0.0/24",
  "10.1.1.0/24",
  "10.1.2.0/24"
]
private_subnet_1_secondary_ranges = [
  "10.2.0.0/24",
  "10.2.1.0/24"
]

# Instance configuration
instance_count = 2
machine_type  = "e2-micro"
network_tags  = ["web", "dev"]

# Load balancer configuration
create_ipv6_address = false

# Resource labels
resource_labels = {
  environment = "dev"
  managed_by  = "terraform"
  project     = "tf-lab3"
}
```

2. Create `prod.tfvars`:
```hcl
project_id   = "iis-tf-dev"
environment  = "prod"
project_name = "tf-lab3"
region       = "us-central1"
zone         = "us-central1-a"

# Network configuration
private_subnet_cidr_blocks = ["10.10.1.0/24", "10.10.2.0/24"]
public_subnet_cidr_blocks  = ["10.10.3.0/24", "10.10.4.0/24"]
private_subnet_0_secondary_ranges = [
  "10.11.0.0/24",
  "10.11.1.0/24",
  "10.11.2.0/24"
]
private_subnet_1_secondary_ranges = [
  "10.12.0.0/24",
  "10.12.1.0/24"
]

# Instance configuration
instance_count = 3  # More instances for production
machine_type  = "e2-small"  # Larger instance type for production
network_tags  = ["web", "prod"]

# Load balancer configuration
create_ipv6_address = true  # Enable IPv6 in production

# Resource labels
resource_labels = {
  environment = "prod"
  managed_by  = "terraform"
  project     = "tf-lab3"
}
```

### 2. Update Resource Names

1. Update the VPC module name in `main.tf`:
```hcl
module "vpc" {
  # ... existing configuration ...
  network_name = "${var.project_name}-${var.environment}"  # Now includes environment
}
```

2. Update the Cloud Router name:
```hcl
module "cloud_router" {
  # ... existing configuration ...
  name = "nat-router-${var.environment}"
}
```

3. Update the load balancer name:
```hcl
module "lb" {
  # ... existing configuration ...
  name = "lb-${var.environment}-${random_string.lb_id.result}"
}
```

### 3. Test the Configuration

Run plan for each environment to verify the changes:

```bash
# Test development configuration
terraform plan -var-file="dev.tfvars"

# Test production configuration
terraform plan -var-file="prod.tfvars"
```

Notice how:
- Resource names include the environment
- Network CIDR ranges are different
- Production has more/larger instances
- Labels reflect the environment

## Next Steps
In the next lab, you will:
1. Set up a GitLab CI/CD pipeline
2. Configure remote state management
3. Create environment-specific backend configurations 