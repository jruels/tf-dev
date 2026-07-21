# Terraform Variables&#x20;

## Overview

In this lab, you will update the existing `main.tf` file to use variables.

## Set the Instance Name with a Variable

### Create the Lab Directory

1. In **Visual Studio Code**, open the working directory created in the previous lab (`YYYYMMDD/terraform`).
2. Right-click in the **Explorer** pane and select **New Folder**.
3. Name the folder `tf-lab2`.
4. Copy `main.tf` from `tf-lab1` to `tf-lab2`:
   - Right-click on `main.tf` in `tf-lab1` and select **Copy**.
   - Navigate to `tf-lab2`, right-click inside the folder, and select **Paste**.

### Define a Variable for the Instance Name

1. Right-click inside `tf-lab2` and select **New File**.
2. Name the file `variables.tf` and open it.
3. Paste the following variable definition:

```hcl
variable "instance_name" {
  description = "Name for GCE instance"
  type        = string
  default     = "lab2-tf-example"
}
```

### Update `main.tf` to Use the Variable

1. Open `main.tf` in `tf-lab2`.
2. Update the `google_compute_instance` resource block to use the new variable:

```hcl
  name = var.instance_name
  labels = {
    name = var.instance_name
  }
```

3. Rename the resource from `lab1-tf-example` to `lab2-tf-example`:

```hcl
resource "google_compute_instance" "lab2-tf-example" {
 
}
```

## Apply the Configuration

1. Open **Integrated Terminal** in `tf-lab2`.

2. Initialize Terraform:
   ```sh
   terraform init
   ```

3. Create an execution plan:
   ```sh
   terraform plan
   ```

4. Apply the configuration:
   ```sh
   terraform apply
   ```

5. If everything looks correct, type `yes` to confirm and apply the configuration.

6. Apply the configuration again, passing the variable via the command line:

   ```sh
   terraform apply -var 'instance_name=SomeOtherName'
   ```

### Note

Variables passed via the command line are not saved, so they must be set each time unless added to a variable file.

## Query Data with Outputs

### Define Output Values

1. Right-click inside `tf-lab2` and select **New File**.
2. Name the file `outputs.tf` and open it.
3. Paste the following output definitions:

```hcl
output "instance_id" {
  description = "ID of the GCE instance"
  value       = google_compute_instance.lab2-tf-example.id
}

output "instance_internal_ip" {
  description = "Internal IP address of GCE instance"
  value       = google_compute_instance.lab2-tf-example.network_interface[0].network_ip
}

output "instance_external_ip" {
  description = "External IP address of GCE instance"
  value       = google_compute_instance.lab2-tf-example.network_interface[0].access_config[0].nat_ip
}
```

### Inspect Output Values

1. Apply the configuration:
   ```sh
   terraform apply
   ```
2. Query the outputs:
   ```sh
   terraform output
   ```
3. Example output:
   ```
   instance_id = "projects/my-project/zones/us-west1-a/instances/lab2-tf-example"
   instance_internal_ip = "10.128.0.2"
   instance_external_ip = "34.82.123.45"
   ```

Terraform outputs are useful for integrating with other infrastructure components or CI/CD pipelines.

## Cleanup

1. Destroy the infrastructure:
   ```sh
   terraform destroy -auto-approve
   ```

## Congratulations!

You have successfully used Terraform variables and outputs. 