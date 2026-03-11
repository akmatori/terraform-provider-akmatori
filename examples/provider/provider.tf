terraform {
  required_providers {
    akmatori = {
      source = "registry.terraform.io/akmatori/akmatori"
    }
  }
}

provider "akmatori" {
  host     = "https://akmatori.example.com"
  username = "admin"
  password = var.akmatori_password
}

variable "akmatori_password" {
  type      = string
  sensitive = true
}
