# Create a basic auth credential
resource "uptime_credential" "api_auth" {
  display_name    = "API Basic Auth"
  credential_type = "BASIC"
  username        = "admin"
  secret = {
    password = "secret-password"
  }
}

# Create an HTTP check
resource "uptime_check_http" "api" {
  name    = "API Health Check"
  address = "https://api.example.com/health"
  # ... other configuration ...
}

# Link the password to the check as a variable
resource "uptime_service_variable" "api_password" {
  service_id    = uptime_check_http.api.id
  credential_id = uptime_credential.api_auth.id
  variable_name = "auth_password"
  property_name = "password"
}
