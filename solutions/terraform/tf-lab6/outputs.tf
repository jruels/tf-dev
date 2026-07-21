output "website_url" {
  description = "URL of the website"
  value       = module.static_website_bucket.url
}

output "bucket_name" {
  description = "Name of the bucket"
  value       = module.static_website_bucket.name
}

output "bucket_self_link" {
  description = "Self link of the bucket"
  value       = module.static_website_bucket.self_link
}