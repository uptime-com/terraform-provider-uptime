# Basic WHOIS/domain expiry check
# expect_string is required and should match text in the domain's WHOIS record
resource "uptime_check_whois" "example" {
  name          = "Domain Expiry Check"
  address       = "example.com"
  expect_string = "IANA"
}

# WHOIS check with expiration threshold
resource "uptime_check_whois" "expiry_alert" {
  name          = "Domain Expiry Alert"
  address       = "google.com"
  expect_string = "Google"
  threshold     = 30 # Alert 30 days before expiry
}

# WHOIS check with full configuration
resource "uptime_check_whois" "full" {
  name           = "Production Domain"
  address        = "github.com"
  expect_string  = "GitHub"
  threshold      = 60
  contact_groups = ["nobody"]
}

# Multiple domain checks
resource "uptime_check_whois" "primary" {
  name          = "Primary Domain"
  address       = "example.com"
  expect_string = "IANA"
  threshold     = 30
}

resource "uptime_check_whois" "secondary" {
  name          = "Secondary Domain"
  address       = "example.net"
  expect_string = "IANA"
  threshold     = 30
}
