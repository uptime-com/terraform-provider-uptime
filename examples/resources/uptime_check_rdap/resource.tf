# Basic RDAP check (domain registration data)
resource "uptime_check_rdap" "example" {
  name    = "Domain RDAP Check"
  address = "example.com"
}

# RDAP check with expiration threshold
resource "uptime_check_rdap" "expiry" {
  name      = "Domain Expiry RDAP"
  address   = "mycompany.com"
  threshold = 30
}

# RDAP check with full configuration
resource "uptime_check_rdap" "full" {
  name           = "Production Domain RDAP"
  address        = "prod.example.com"
  threshold      = 60
  contact_groups = ["nobody"]
}
