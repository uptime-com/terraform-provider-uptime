# Retrieve all private monitoring locations
data "uptime_private_locations" "all" {}

# Output all private locations
output "all_private_locations" {
  value       = data.uptime_private_locations.all.locations
  description = "All private monitoring locations"
}

# Use private locations in a check. Check locations are matched by the
# "location" attribute, not the user-defined "name".
locals {
  private_locations = [
    for loc in data.uptime_private_locations.all.locations : loc.location
  ]
}

resource "uptime_check_http" "private" {
  name      = "Private Location Check"
  address   = "https://api.internal.example.com"
  locations = local.private_locations
}
