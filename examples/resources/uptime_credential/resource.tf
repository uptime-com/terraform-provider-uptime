# Basic authentication credential
resource "uptime_credential" "example" {
  display_name    = "API Credentials"
  credential_type = "BASIC"
  secret = {
    password = "secure-password"
  }
}

# Basic auth with username
resource "uptime_credential" "basic_auth" {
  display_name    = "Admin Credentials"
  credential_type = "BASIC"
  description     = "Credentials for admin panel access"
  username        = "admin"
  secret = {
    password = "admin-password"
  }
}

# Token credential
resource "uptime_credential" "token" {
  display_name    = "API Token"
  credential_type = "TOKEN"
  description     = "Bearer token for API authentication"
  secret = {
    secret = "your-api-token-here"
  }
}

# Certificate credential
# Note: Use file() function to read certificate files
resource "uptime_credential" "certificate" {
  display_name    = "Client Certificate"
  credential_type = "CERTIFICATE"
  description     = "TLS client certificate for mTLS"
  secret = {
    certificate = "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
    key         = "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
    passphrase  = "key-passphrase"
  }
}
