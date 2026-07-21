variable "instance_count" {
  description = "Number of instances to create"
  type        = number
}

variable "name_prefix" {
  description = "Prefix for instance names"
  type        = string
}

variable "machine_type" {
  description = "Machine type for instances"
  type        = string
  default     = "e2-micro"
}

variable "zone" {
  description = "Zone where instances will be created"
  type        = string
}

variable "network" {
  description = "Network for instances"
  type        = string
}

variable "subnetwork" {
  description = "Subnetwork for instances"
  type        = string
}

variable "enable_public_ip" {
  description = "Whether to enable public IP for instances"
  type        = bool
  default     = false
}

variable "network_tags" {
  description = "Network tags for instances"
  type        = list(string)
  default     = []
}

variable "labels" {
  description = "Labels to apply to instances"
  type        = map(string)
  default     = {}
} 