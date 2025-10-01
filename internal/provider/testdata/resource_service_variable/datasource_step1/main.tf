variable "credential_name" {
  type = string
}

variable "password" {
  type = string
}

resource "uptime_credential" "test" {
  display_name    = var.credential_name
  credential_type = "BASIC"
  username        = "testuser"
  secret = {
    password = var.password
  }
}
