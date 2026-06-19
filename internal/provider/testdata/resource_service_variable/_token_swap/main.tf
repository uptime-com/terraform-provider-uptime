variable "credential_name_a" {
  type = string
}

variable "credential_name_b" {
  type = string
}

variable "token" {
  type = string
}

variable "use_b" {
  type = bool
}

resource "uptime_credential" "a" {
  display_name    = var.credential_name_a
  credential_type = "TOKEN"
  secret = {
    secret = var.token
  }
}

resource "uptime_credential" "b" {
  display_name    = var.credential_name_b
  credential_type = "TOKEN"
  secret = {
    secret = var.token
  }
}

resource "uptime_check_http" "test" {
  name    = "test-check-for-service-variable-token"
  address = "https://example.com"
}

resource "uptime_service_variable" "test" {
  service_id    = uptime_check_http.test.id
  credential_id = var.use_b ? uptime_credential.b.id : uptime_credential.a.id
  variable_name = "token_raven_token"
  property_name = "secret"
}
