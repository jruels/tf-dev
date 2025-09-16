resource "google_storage_bucket" "website" {
  name          = var.bucket_name
  location      = var.location
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD", "OPTIONS"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  labels = var.labels
}

resource "google_storage_bucket_iam_member" "admin" {
  bucket = google_storage_bucket.website.name
  role   = "roles/storage.objectAdmin"
  member = "user:${var.admin_email}"
}