# Retrieve all available monitoring locations
data "uptime_locations" "all" {}

# Output all available locations
output "all_locations" {
  value       = data.uptime_locations.all.locations
  description = "All available monitoring locations"
}

# Filter for US locations only
locals {
  us_locations = [
    for loc in data.uptime_locations.all.locations :
    loc.name if can(regex("^US", loc.name))
  ]
}

output "us_locations" {
  value       = local.us_locations
  description = "US-based monitoring locations"
}

# Filter for European locations
locals {
  eu_locations = [
    for loc in data.uptime_locations.all.locations :
    loc.name if can(regex("^EU", loc.name))
  ]
}

output "eu_locations" {
  value       = local.eu_locations
  description = "European monitoring locations"
}

# Use locations data in a check
resource "uptime_check_http" "multi_region" {
  name      = "Multi-Region Check"
  address   = "https://api.example.com"
  locations = local.us_locations
}
