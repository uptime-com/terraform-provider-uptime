# Retrieve all credentials
data "uptime_credentials" "all" {}

# Filter for a specific credential
locals {
  basic_auth_cred = [
    for cred in data.uptime_credentials.all.credentials :
    cred if cred.credential_type == "BASIC" && cred.display_name == "My Basic Auth"
  ][0]
}

# Create a check
resource "uptime_check_http" "example" {
  name    = "Example Check"
  address = "https://example.com"
  # ... other configuration ...
}

# Link the credential to the check using a service variable
resource "uptime_service_variable" "auth" {
  service_id    = uptime_check_http.example.id
  credential_id = local.basic_auth_cred.id
  variable_name = "auth_password"
  property_name = "password"
}
