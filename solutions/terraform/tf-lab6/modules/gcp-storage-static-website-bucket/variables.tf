variable "project_id" {
  description = "The ID of the project where the bucket will be created"
  type        = string
}

variable "bucket_name" {
  description = "Name of storage bucket. Must be unique"
  type        = string
}

variable "location" {
  description = "Location for the storage bucket"
  type        = string
  default     = "US"
}

variable "admin_email" {
  description = "Email address of the bucket admin"
  type        = string
}

variable "labels" {
  description = "Labels to set on the bucket"
  type        = map(string)
  default     = {}
}