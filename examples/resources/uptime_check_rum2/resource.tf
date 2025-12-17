# Basic RUM (Real User Monitoring) check
resource "uptime_check_rum2" "example" {
  name    = "Website RUM"
  address = "www.example.com"
}

# RUM check for SPA application
resource "uptime_check_rum2" "spa" {
  name    = "SPA Real User Monitoring"
  address = "app.example.com"
}

# RUM check with full configuration
resource "uptime_check_rum2" "full" {
  name           = "Production RUM"
  address        = "www.example.com"
  contact_groups = ["nobody"]
}
