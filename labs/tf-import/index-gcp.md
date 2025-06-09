# Terraform - import existing resources

## Overview
In this lab, you will create Google Cloud resources using the console and import them into Terraform management. 

## Create instances in Google Cloud Console
Create three Compute Engine instances in the Google Cloud Console. 

1. Navigate to the Google Cloud Console and use the search bar at the top to search for "Compute Engine" and click on it.

2. In the Compute Engine dashboard, click on "CREATE INSTANCE".

3. For each instance, configure the following:
   - Name: Leave as default (will be changed later)
   - Region/Zone: Select your preferred region (e.g., us-central1-a)
   - Machine configuration: Select "E2" series and "e2-micro"
   - Leave other settings as default

4. Click "Create" and repeat this process two more times for a total of 3 instances.

## Create Terraform configuration 
While waiting for the instances to launch, create a new working directory and configuration file. 

1. In **Visual Studio Code**, open the working directory created in the previous lab (`YYYYMMDD/terraform`).

2. Right-click in the **Explorer** pane and select **New Folder**.

3. Name the folder `tf-lab4`.

4. Right-click `tf-lab4` and select **New File**.

5. Name the file `main.tf` and add a resource with the following attributes:

   - type: `google_compute_instance`
   - name: `tf-example-import`
   - machine_type: `e2-micro`
   - zone: Your selected zone (e.g., us-central1-a)
   - count: `3`
   - boot_disk:
     - initialize_params:
       - image: `debian-cloud/debian-11`
   - network_interface:
     - network: `default`
   - labels:
     - role: `terraform`
     - name: `tf-example-import-${count.index}`

Remember, this resource block is for three instances. You will need to add the `count.index` to the instance name and labels. If you get stuck, ask the instructor for assistance.

## Import the configuration 
Now that you've created the instances and the Terraform configuration, first initialize Terraform:

```bash
terraform init
```

Then use the `terraform import` command to import the existing instances. 

The import command format for GCP Compute Engine instances is:
```
terraform import 'google_compute_instance.tf-example-import[0]' PROJECT_ID/ZONE/INSTANCE_NAME
```

You'll need to run this command for each instance, incrementing the index [0], [1], [2] and using the corresponding instance names from the console.

If you get stuck, check the help page `terraform import --help` or the [Terraform GCP Provider documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance#import).

## Apply Changes

After importing, if you run:
```bash
terraform apply
```

Note that this will rename your instances to match the names in your Terraform configuration (tf-example-import-0, tf-example-import-1, tf-example-import-2). This is because the configuration specifies these names, and Terraform will modify the instances to match the configuration.

## Cleanup

To remove all resources:
```bash
terraform destroy
```

This will delete all three instances that were imported.

# Congrats

