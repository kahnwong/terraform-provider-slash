terraform {
  required_providers {
    slash = {
      source = "hashicorp.com/kahnwong/slash"
    }
  }
}

provider "slash" {}

# data "hashicups_coffees" "example" {}

resource "slash_shortcut" "example" {
  name  = "mbs"
  link  = "https://microbin.karnwong.me"
  title = "Microbin"
}
