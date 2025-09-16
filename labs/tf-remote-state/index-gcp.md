# Terraform - remote state configuration

## Overview 
In this lab, you will create a Google Cloud Storage bucket and migrate the Terraform state to a remote backend. 

### Manual Bucket Creation

First, manually create a GCS bucket for remote state:

1. Go to the Google Cloud Console and select your project from the project dropdown at the top of the page
2. Navigate to Cloud Storage > Buckets
3. Click "CREATE BUCKET"
4. Enter a unique name for your bucket (e.g., "remote-state-YOUR_NAME")
5. Choose "US" for Location type
6. Leave other settings as default
7. Under "Labels" add:
   - Key: name
   - Value: remote-state-backend
8. Click "CREATE"

## Migrate the state
Now that we've created a Cloud Storage bucket, we need to migrate the state to the remote backend. 

When creating the backend configuration remember to replace the `bucket` with the name of the bucket you created. 

Create `backend.tf` with the following:
```hcl
terraform {
  backend "gcs" {
    bucket = "remote-state-jrs"
    prefix = "terraform/state"
  }
}
```

## Reinitialize Terraform 
Now that you have created the GCS bucket and configured the `backend.tf` you must run `terraform init` to migrate the state to the new remote backend. 

If prompted to migrate the existing state type `yes`

If everything is successful, you should see a message that the backend was migrated. 

## Cleanup

```bash
terraform destroy
```

## Congratulations

You have successfully created a Cloud Storage bucket and migrated to a remote backend state. 