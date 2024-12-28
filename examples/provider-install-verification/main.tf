terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/kahnwong/slash"
    }
  }
}

provider "hashicups" {}

# data "hashicups_coffees" "example" {}
