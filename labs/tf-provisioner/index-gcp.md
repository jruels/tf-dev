# Terraform - provisioners

## Overview 
In this lab, you will update `main.tf` to include provisioners.

Provisioners allow you to run shell scripts on the local machine or remote resources. You can also use vendor provisioners from Chef, Puppet, and Salt.

## Add provisioners
This lab updates the `main.tf` in the `tf-lab2` directory. 

Add a `local-exec` provisioner to your Google Compute Engine instance with the following attributes: 
- command: Echo the external IP addresses into a file named `external_ips.txt`
```hcl
provisioner "local-exec" {
  command = "echo ${self.network_interface[0].access_config[0].nat_ip} >> external_ips.txt"
}
```

Add another `local-exec` provisioner with the following attributes: 
- command: Echo the internal IP addresses into a file named `internal_ips.txt`
```hcl
provisioner "local-exec" {
  command = "echo ${self.network_interface[0].network_ip} >> internal_ips.txt"
}
```

Note: In GCP, we refer to:
- Public IPs as "external IPs" or "NAT IPs"
- Private IPs as "internal IPs" or "network IPs"

The IP addresses in GCP are accessed differently from AWS:
- External IP: `network_interface[0].access_config[0].nat_ip`
- Internal IP: `network_interface[0].network_ip`

## Apply and Check Results

Run terraform apply to create the instance and trigger the provisioners:
```bash
terraform apply -auto-approve
```

Once complete, you can view the IP addresses that were captured:
```bash
cat external_ips.txt  # View external (public) IP
cat internal_ips.txt  # View internal (private) IP
```

## Cleanup

Run the following to clean up the resources

```bash
terraform destroy -auto-approve
```

## Congratulations

You have successfully added provisioners to your GCP Compute Engine instance to capture both external and internal IP addresses. 