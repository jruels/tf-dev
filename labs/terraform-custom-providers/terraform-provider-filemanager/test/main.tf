terraform {
  required_providers {
    filemanager = {
      source  = "registry.terraform.io/donis/filemanager"
      version = "0.1.0"
    }
  }
}

resource "filemanager_file" "example_file" {
  path = "c:/code/test.txt"
}

output "is_large_file" {
  value = filemanager_file.example_file.is_large
}

resource "filemanager_directory" "example_dir" {
  path = "c:/windows/system32"
}

output "large_files_in_directory" {
  value = filemanager_directory.example_dir.large_files
}