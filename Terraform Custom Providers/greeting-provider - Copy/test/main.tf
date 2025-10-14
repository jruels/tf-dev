terraform {
  required_providers {
    greeting2 = {
      source  = "donis/greeting2"
      version = "1.0.0"
    }
  }
}

provider "greeting2" {}

resource "greeting_message" "hello" {}