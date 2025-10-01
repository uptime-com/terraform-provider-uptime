variable "credential_name" {
  type = string
}

variable "password" {
  type = string
}

variable "variable_name" {
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

data "uptime_credentials" "all" {}

locals {
  # Filter credentials datasource to find our test credential by display_name
  test_credential = [
    for cred in data.uptime_credentials.all.credentials :
    cred if cred.display_name == var.credential_name
  ][0]
}

resource "uptime_check_http" "test" {
  name    = "test-check-for-service-variable"
  address = "https://example.com"
}

resource "uptime_service_variable" "test" {
  service_id    = uptime_check_http.test.id
  credential_id = local.test_credential.id
  variable_name = var.variable_name
  property_name = "password"
}
