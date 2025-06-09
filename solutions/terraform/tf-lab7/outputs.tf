output "website_url" {
  description = "URL of the website"
  value       = "https://storage.googleapis.com/${google_storage_bucket.website.name}/index.html"
}