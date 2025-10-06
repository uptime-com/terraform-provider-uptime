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

resource "uptime_check_http" "test" {
  name    = "test-check-for-service-variable"
  address = "https://example.com"
}

resource "uptime_service_variable" "test" {
  service_id    = uptime_check_http.test.id
  credential_id = uptime_credential.test.id
  variable_name = var.variable_name
  property_name = "password"
}
