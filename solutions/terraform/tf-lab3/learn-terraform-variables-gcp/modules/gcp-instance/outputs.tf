output "instance_ids" {
  description = "IDs of created instances"
  value       = google_compute_instance.vm[*].id
}

output "instance_names" {
  description = "Names of created instances"
  value       = google_compute_instance.vm[*].name
}

output "instance_self_links" {
  description = "Self-links of created instances"
  value       = google_compute_instance.vm[*].self_link
}

output "internal_ips" {
  description = "Internal IPs of created instances"
  value       = google_compute_instance.vm[*].network_interface[0].network_ip
}

output "external_ips" {
  description = "External IPs of created instances (if enabled)"
  value       = [for instance in google_compute_instance.vm[*] : instance.network_interface[0].access_config[0].nat_ip if length(instance.network_interface[0].access_config) > 0]
} 