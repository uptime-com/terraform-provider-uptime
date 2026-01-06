# Basic SSL certificate check
resource "uptime_check_sslcert" "example" {
  name    = "SSL Certificate Check"
  address = "example.com"
}

# SSL check with expiration threshold
resource "uptime_check_sslcert" "expiry_alert" {
  name      = "SSL Expiry Alert"
  address   = "www.example.com"
  threshold = 30 # Alert 30 days before expiry
}

# SSL check with full configuration
resource "uptime_check_sslcert" "full" {
  name           = "Production SSL"
  address        = "secure.example.com"
  threshold      = 14
  contact_groups = ["nobody"]
}

# SSL check for multiple domains
resource "uptime_check_sslcert" "api" {
  name      = "API SSL Certificate"
  address   = "api.example.com"
  threshold = 30
}

resource "uptime_check_sslcert" "cdn" {
  name      = "CDN SSL Certificate"
  address   = "cdn.example.com"
  threshold = 30
}
