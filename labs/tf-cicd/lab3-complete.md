# GitLab CI/CD Lab: Complete Guide

This lab converts the GitHub Actions AWS lab to use GitLab CI/CD and Google Cloud Platform (GCP). You'll set up a complete CI/CD pipeline that manages GCP infrastructure using Terraform.

## Part 1: Initial Setup and Configuration

### Generate GCP Service Account
1. Create a new project in GCP Console (if not already done)
2. Enable required APIs:
   ```bash
   gcloud services enable compute.googleapis.com
   gcloud services enable iam.googleapis.com
   gcloud services enable cloudresourcemanager.googleapis.com
   ```
3. Create service account and download key (from Lab 1.5)

### Create GitLab Repository
1. Log into GitLab
2. Create new project: "terraform-gcp-lab"
3. Clone the repository:
   ```bash
   git clone git@gitlab.com:your-username/terraform-gcp-lab.git
   cd terraform-gcp-lab
   ```

## Part 2: GitLab and GCP Configuration

### Add GitLab CI/CD Variables
1. Go to Settings → CI/CD → Variables
2. Add the following variables:
   ```
   Name: PROJECT_ID
   Value: your-gcp-project-id
   Type: Variable
   Protect: No
   Mask: No

   Name: GCP_SA_KEY
   Value: (contents of your service account key JSON)
   Type: File
   Protect: Yes
   Mask: Yes
   ```

### Create GitLab Environment
1. Go to Settings → CI/CD → Environments
2. Click "New environment"
3. Name it: production
4. Add deployment freeze periods if desired
5. Save the environment

## Part 3: GCP Configuration

### Create GCS Backend
1. Create a GCS bucket:
   ```bash
   export PROJECT_ID="your-project-id"
   gsutil mb -l us-central1 gs://${PROJECT_ID}-terraform-state
   gsutil versioning set on gs://${PROJECT_ID}-terraform-state
   ```

## Part 4: Repository Configuration

### Set Up Infrastructure Code

1. Create backend configuration (backend.tf):
```hcl
terraform {
  backend "gcs" {
    bucket = "your-project-id-terraform-state"
    prefix = "dev/terraform.tfstate"
  }
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
```

2. Create variables (variables.tf):
```hcl
variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "GCP Zone"
  type        = string
  default     = "us-central1-a"
}

variable "instance_type" {
  description = "Instance type"
  type        = string
  default     = "e2-micro"
}
```

3. Create network configuration (network.tf):
```hcl
resource "google_compute_network" "vpc" {
  name                    = "terraform-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "terraform-subnet"
  ip_cidr_range = "10.0.1.0/24"
  network       = google_compute_network.vpc.id
  region        = var.region
}

resource "google_compute_firewall" "web" {
  name    = "web-access"
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["80", "443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["web"]
}
```

4. Create compute configuration (compute.tf):
```hcl
resource "google_compute_instance" "web" {
  name         = "web-server"
  machine_type = var.instance_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnet.id
    access_config {
      // Ephemeral public IP
    }
  }

  metadata_startup_script = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y apache2
    echo "Welcome to the GitLab CI/CD & Terraform Lab!!" > /var/www/html/index.html
    systemctl start apache2
    systemctl enable apache2
  EOF

  tags = ["web"]
}

output "web_ip" {
  value = google_compute_instance.web.network_interface[0].access_config[0].nat_ip
}
```

### Create GitLab CI/CD Pipeline

Create `.gitlab-ci.yml`:
```yaml
image: hashicorp/terraform:1.5

variables:
  TF_VAR_project_id: ${PROJECT_ID}
  GOOGLE_APPLICATION_CREDENTIALS: ${CI_PROJECT_DIR}/credentials.json

stages:
  - validate
  - plan
  - apply
  - destroy

before_script:
  - echo $GCP_SA_KEY > ${CI_PROJECT_DIR}/credentials.json

validate:
  stage: validate
  script:
    - terraform init
    - terraform fmt -check
    - terraform validate
  rules:
    - if: $CI_MERGE_REQUEST_IID
    - if: $CI_COMMIT_BRANCH == "main"

plan:
  stage: plan
  script:
    - terraform init
    - terraform plan -out=plan.tfplan
  artifacts:
    paths:
      - plan.tfplan
    expire_in: 1 week
  rules:
    - if: $CI_MERGE_REQUEST_IID
    - if: $CI_COMMIT_BRANCH == "main"

apply:
  stage: apply
  script:
    - terraform init
    - terraform apply plan.tfplan
  dependencies:
    - plan
  environment:
    name: production
    url: http://$(<terraform output -raw web_ip)
  when: manual
  rules:
    - if: $CI_COMMIT_BRANCH == "main"

destroy:
  stage: destroy
  script:
    - terraform init
    - terraform destroy -auto-approve
  environment:
    name: production
    action: stop
  when: manual
  rules:
    - if: $CI_MERGE_REQUEST_IID
```

### Initial Deployment

1. Create and push feature branch:
```bash
git checkout -b feature/initial-setup
git add .
git commit -m "Initial infrastructure setup"
git push -u origin feature/initial-setup
```

2. Create Merge Request:
   - Go to GitLab repository
   - Click "Merge Requests" → "New merge request"
   - Source branch: feature/initial-setup
   - Target branch: main
   - Create merge request

3. Watch the Pipeline:
   - Pipeline will automatically run validate and plan
   - Review the plan output carefully
   - If everything looks correct, merge the request
   - Go to the pipeline in main branch
   - Manually approve the apply job

4. Verify Deployment:
   - After successful apply, get the web server IP from the pipeline output
   - Open the IP in your browser
   - You should see: "Welcome to the GitLab CI/CD & Terraform Lab!!"

## Part 5: Update Instance Type

1. Create new feature branch:
```bash
git checkout -b feature/upgrade-instance
```

2. Update variables.tf:
```hcl
variable "instance_type" {
  description = "Instance type"
  type        = string
  default     = "e2-small"  # Changed from e2-micro
}
```

3. Commit and push:
```bash
git add .
git commit -m "Upgrade instance type to e2-small"
git push -u origin feature/upgrade-instance
```

4. Create and process merge request:
   - Create new merge request
   - Review the plan (should show 1 change)
   - Merge the request
   - Approve the apply job in the main branch pipeline

5. Verify:
   - Check GCP Console that instance type changed
   - Verify website still works

## Part 6: Cleanup

1. To destroy all resources:
   - Go to main branch pipeline
   - Run the destroy job manually
   - Verify in GCP Console that all resources are removed

## Lab Summary

In this lab you:
- Set up GitLab CI/CD with GCP
- Created a complete infrastructure pipeline
- Managed infrastructure changes through merge requests
- Implemented secure credential handling
- Created a production environment with manual approvals
- Learned to manage state in GCS

## Troubleshooting Guide

### Common Issues and Solutions

#### Pipeline Authentication Issues
- "Could not load credentials" error:
  - Verify GCP_SA_KEY is properly formatted
  - Check that the service account has required permissions

#### GCS Backend Issues
- "Access denied" for state bucket:
  - Verify service account has Storage Admin role
  - Check bucket name matches exactly in backend.tf

#### Environment Issues
- Manual approval not appearing:
  - Verify environment name matches exactly in .gitlab-ci.yml
  - Check that environment is created in GitLab settings 