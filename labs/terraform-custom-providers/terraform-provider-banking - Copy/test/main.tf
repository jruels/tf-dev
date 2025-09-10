terraform {
  required_providers {
    banking = {
      source = "registry.terraform.io/example/banking"
    }
  }
}

provider "banking" {
  db_host     = "localhost"
  db_port     = 5432
  db_user     = "postgres"
  db_password = "Post1260"
  db_name     = "bankingdb"
}

resource "banking_customer_account" "customer1" {
  first_name   = "Alice"
  last_name    = "Doe"
  email        = "alice@example.com"
  account_type = "savings"
  balance      = 1500.75
}
