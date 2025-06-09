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