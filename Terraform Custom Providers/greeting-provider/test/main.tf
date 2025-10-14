terraform {
  required_providers {
    greeting = {
      source  = "registry.terraform.io/donis/greeting"
      version = "1.0.0"
    }
  }
}

provider "greeting" {}

resource "greeting_message" "hello" {}

output "greeting_output" {
  value = greeting_message.hello.message
}
