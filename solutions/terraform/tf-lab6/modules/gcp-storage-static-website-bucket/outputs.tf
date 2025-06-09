output "url" {
  description = "The URL of the website"
  value       = "https://storage.googleapis.com/${google_storage_bucket.website.name}/index.html"
}

output "name" {
  description = "Name of the bucket"
  value       = google_storage_bucket.website.name
}

output "self_link" {
  description = "Self link of the bucket"
  value       = google_storage_bucket.website.self_link
}