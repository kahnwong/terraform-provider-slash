terraform {
  required_providers {
    slash = {
      source = "hashicorp.com/kahnwong/slash"
    }
  }
}

provider "slash" {
  host         = "http://localhost:5231"
  access_token = "foobarbaz"
}

resource "slash_shortcut" "example" {
  name  = "mb"
  link  = "https://foo.bar"
  title = "Microbin"
}