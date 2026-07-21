output "network_name" {
  description = "The name of the VPC being created"
  value       = module.vpc.network_name
}

output "subnets_names" {
  description = "The names of the subnets being created"
  value       = module.vpc.subnets_names
}