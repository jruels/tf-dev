# Terraform - working with variables

## Overview 
Terraform supports several variable types, including bool, string, and numbers.

The following steps continue where the previous lab left off. The changes below should be applied to the `tf-lab3` directory.

## IPv6 Support

Use a bool type variable to control whether your load balancer automatically gets an IPv6 address. Add a declaration for create_ipv6_address to variables.tf:

- variable name: create_ipv6_address
- description: Allocate a new IPv6 address for the load balancer. Conflicts with manually specified ipv6_address.
- type: bool
- default: false

Find the load balancer module block in main.tf (module "lb") and add create_ipv6_address = var.create_ipv6_address to its configuration. This setting should be at the same level as the project, name, and target_tags arguments.

Leave the ipv6_address setting as undefined. This allows users to either let GCP automatically allocate an address (using create_ipv6_address) or specify one manually (using ipv6_address) if needed.

Try applying with IPv6 enabled:
```bash
terraform apply -var="create_ipv6_address=true"
```

Then try with IPv6 disabled:
```bash
terraform apply -var="create_ipv6_address=false"
```

When you write Terraform modules you intend to re-use, you will usually want to make as many attributes configurable with variables as possible, to make your module more flexible for use in more situations.

When you write Terraform configuration for a specific project, you may choose to leave some attributes with hard-coded values when it doesn't make sense to allow users to configure them.

## List public and private subnets
So far you have used "simple" variables. These variables all have single values. Now let's use some more complex variables. Terraform supports collection variable types that contain more than one value. Terraform supports several collection variables types. 
- List: A sequence of values of the same type.
- Map: A lookup table, matching keys to values, all of the same type.
- Set: An unordered collection of unique value, all of the same type.

The following lab will use lists and a map, which are the most commonly used of these types. Sets are useful when a unique collection of values is needed, and the order of the items in the collection does not matter.

A likely place to use list variables is when setting the `secondary_ip_ranges` for subnets. Make this configuration easier to use while still being customizable by using lists along with the `slice()` function.

Add the following variable declaration to `variables.tf`.

```hcl
variable "private_subnet_0_secondary_ranges" {
  description = "Available CIDR blocks for secondary IP ranges."
  type        = list(string)
  default     = [
    "192.168.10.0/24",
    "192.168.20.0/24",
    "192.168.30.0/24"
  ]
}

variable "private_subnet_1_secondary_ranges" {
  description = "Available CIDR blocks for secondary IP ranges."
  type        = list(string)
  default     = [
    "192.168.40.0/24",
    "192.168.50.0/24",
    "192.168.60.0/24"
  ]
}
```

Notice that the type for the list variables is `list(string)`. Each element in these lists must be a string. List elements must all be the same type, but can be any type, including complex types like `list(list)` and `list(map)`.

Like lists and arrays used in most programming languages, you can refer to individual items in a list by index, starting with 0. Terraform also includes several functions that allow you to manipulate lists and other variable types.

Use the `slice()` function to get a subset of these lists.

The Terraform console command opens an interactive console that you can use to evaluate expressions in the context of your configuration. This can be very useful when working with and troubleshooting variable definitions.

Initialize the directory

```
terraform init
```

Open a console with the `terraform console` command.

```sh
terraform console
>
```

Now inspect the list of secondary IP ranges in the Terraform console.

Refer to the variable by name to return the entire list.

```sh
var.private_subnet_0_secondary_ranges
```

output: 
```sh
tolist([
  "192.168.10.0/24",
  "192.168.20.0/24",
  "192.168.30.0/24"
])
```

Now use the following to retrieve the second element from the list: 
```sh
var.private_subnet_0_secondary_ranges[1]
```

output: 
```sh
"192.168.20.0/24"
```

Now use the `slice()` function to return the first two elements from the list.
```sh
slice(var.private_subnet_0_secondary_ranges, 0, 2)
```

output:
```sh
tolist([
  "192.168.10.0/24",
  "192.168.20.0/24"
])
```

The `slice()` function takes three arguments: the list to slice, the index to start from, and the number of elements. It returns a new list with the specified elements copied ("sliced") from the original list.

Leave the console by typing `exit` or pressing `Control-D`.

Now that we understand how the `slice()` function works it's time to use it in our configuration. 

In the `main.tf` use the slice function to extract a subnet of the CIDR block lists when defining your VPC's subnet configuration.

Remove lines with `-` and add lines with `+`

```hcl
module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "~> 7.0"

  project_id   = var.project_id
  network_name = "${var.project_name}-${var.environment}"
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name   = "private-subnet-0"
      subnet_ip     = var.private_subnet_cidr_blocks[0]
      subnet_region = var.region
      secondary_ip_ranges = [
        {
          range_name    = "secondary-range-0"
          ip_cidr_range = slice(var.private_subnet_0_secondary_ranges, 0, 3)[0]
        },
        {
          range_name    = "secondary-range-1"
          ip_cidr_range = slice(var.private_subnet_0_secondary_ranges, 0, 3)[1]
        },
        {
          range_name    = "secondary-range-2"
          ip_cidr_range = slice(var.private_subnet_0_secondary_ranges, 0, 3)[2]
        }
      ]
    },
    {
      subnet_name   = "private-subnet-1"
      subnet_ip     = var.private_subnet_cidr_blocks[1]
      subnet_region = var.region
      secondary_ip_ranges = [
        {
          range_name    = "subnet-1-secondary-0"
          ip_cidr_range = var.private_subnet_1_secondary_ranges[0]
        },
        {
          range_name    = "subnet-1-secondary-1"
          ip_cidr_range = var.private_subnet_1_secondary_ranges[1]
        }
      ]
    }
  ]
}
```

This way, users of this configuration can specify the number of subnets and secondary IP ranges they want without worrying about defining CIDR blocks. The slice function is used to extract the appropriate ranges from the list of available CIDR blocks.

## Map resource labels

Each of the resources declared in `main.tf` includes two labels: `project_name` and `environment`. Assign these labels with a map variable type.

Declare a new map variable for resource labels in `variables.tf`.

```hcl
variable "resource_labels" {
  description = "Labels to set for all resources"
  type        = map(string)
  default     = {
    project     = "project-alpha",
    environment = "dev"
  }
}
```

Setting the type to `map(string)` tells Terraform to expect strings for the values in the map. Map keys are always strings. Like dictionaries or maps from programming languages, you can retrieve values from a map with the corresponding key. See how this works with the Terraform console.

```sh
var.resource_labels["environment"]
```

output: 
```sh
"dev"
```

Exit the console 

Now, replace the hard-coded labels in `main.tf` with the new variable. 

Remove lines with `-` and add lines with `+`

```hcl
-  labels = {
-    project     = "project-alpha",
-    environment = "dev"
-  }
+  labels = var.resource_labels

# ... replace all occurrences of `labels = {...}`
```

The hard-coded labels are used multiple times in this configuration, be sure to replace them all.

Apply these changes.

The value of `project` label has changed, so you will be prompted to apply the changes. Respond with `yes` to confirm the changes.

## Cleanup
Run `terraform destroy` to remove resources.

# Congrats 