# Retrieve all private monitoring locations
data "uptime_private_locations" "all" {}

# Output all private locations
output "private_locations" {
  value       = data.uptime_private_locations.all.locations
  description = "All private monitoring locations"
}

# Use private locations in a check
resource "uptime_check_http" "internal_service" {
  name      = "Internal Service Check"
  address   = "https://internal.example.com"
  locations = [for loc in data.uptime_private_locations.all.locations : loc.location]
}
