# Terraform CI/CD Lab 1.5: GitLab and GCP Setup

## Overview
Before we can implement CI/CD pipelines, we need to set up our GitLab environment and GCP service account. In this lab, you will:
1. Create a GitLab account and repository
2. Create a GCP service account with appropriate permissions
3. Configure GitLab with GCP credentials

## Steps

### 1. Set Up GitLab Account and Repository

1. Create a GitLab Account:
   - Go to [GitLab.com](https://gitlab.com)
   - Click "Register now"
   - Fill in your details and create your account
   - Verify your email address

2. Create a New Project:
   - Click "Create new project"
   - Select "Create blank project"
   - Project name: `terraform-labs`
   - Visibility Level: Private
   - Click "Create project"

3. Configure Git Locally:
```bash
# Configure Git if you haven't already
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# Clone the new repository
git clone git@gitlab.com:your-username/terraform-labs.git
cd terraform-labs

# Copy your existing tf-lab3 files
cp -r /path/to/tf-lab3 .
```

### 2. Create GCP Service Account

1. Create the service account:
```bash
# Set your project ID
export PROJECT_ID="your-project-id"

# Create service account
gcloud iam service-accounts create gitlab-terraform \
    --display-name="GitLab Terraform Pipeline SA" \
    --project="${PROJECT_ID}"
```

2. Grant necessary permissions:
```bash
# Get the service account email
export SA_EMAIL="gitlab-terraform@${PROJECT_ID}.iam.gserviceaccount.com"

# Grant required roles
gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="roles/compute.admin"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="roles/storage.admin"

gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="roles/iam.serviceAccountUser"
```

3. Create and download service account key:
```bash
gcloud iam service-accounts keys create gitlab-terraform-key.json \
    --iam-account="${SA_EMAIL}" \
    --project="${PROJECT_ID}"
```

### 3. Configure GitLab CI/CD Variables

1. In GitLab, navigate to:
   - Settings > CI/CD > Variables
   - Click "Add variable"

2. Add the following variables:
   - `PROJECT_ID`:
     - Key: PROJECT_ID
     - Value: your-project-id
     - Type: Variable
     - Environment scope: All (default)
     - Protect variable: No
     - Mask variable: No

   - `GCP_SA_KEY`:
     - Key: GCP_SA_KEY
     - Value: (contents of gitlab-terraform-key.json)
     - Type: File
     - Environment scope: All (default)
     - Protect variable: Yes
     - Mask variable: Yes

### 4. Test the Setup

1. Create a simple test file to verify the setup:
```bash
# Create test.tf
cat > test.tf << EOF
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
  region  = "us-central1"
}

variable "project_id" {}

output "project_number" {
  value = data.google_project.current.number
}

data "google_project" "current" {
  project_id = var.project_id
}
EOF
```

2. Create a basic pipeline file:
```bash
# Create .gitlab-ci.yml
cat > .gitlab-ci.yml << EOF
image: hashicorp/terraform:1.5

variables:
  TF_VAR_project_id: \${PROJECT_ID}
  GOOGLE_APPLICATION_CREDENTIALS: \${CI_PROJECT_DIR}/credentials.json

test_setup:
  script:
    - echo \$GCP_SA_KEY > \${CI_PROJECT_DIR}/credentials.json
    - terraform init
    - terraform plan
EOF
```

3. Commit and push the test files:
```bash
git add .
git commit -m "Add test configuration"
git push origin main
```

4. Verify the pipeline:
   - Go to GitLab > CI/CD > Pipelines
   - Check that the pipeline runs successfully
   - If there are any errors, review the job logs for troubleshooting

### 5. Clean Up Test Files

Once everything is working:
```bash
git rm test.tf
git commit -m "Remove test configuration"
git push origin main
```

## Important Notes

1. Security Best Practices:
   - Never commit service account keys to Git
   - Use masked and protected variables for sensitive data
   - Follow the principle of least privilege for service account permissions

2. GitLab Settings:
   - Keep your repository private
   - Enable branch protection for `main`
   - Consider enabling merge request approvals

3. GCP Best Practices:
   - Use separate service accounts for different environments
   - Regularly rotate service account keys
   - Monitor service account usage

## Next Steps
Now that you have:
- A working GitLab repository
- GCP service account with appropriate permissions
- Verified CI/CD pipeline setup

You're ready to proceed to Lab 2, where you'll implement the full CI/CD pipeline with state management. 