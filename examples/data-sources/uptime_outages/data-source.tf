# Retrieve outage history from your Uptime.com account
data "uptime_outages" "all" {}

# Filter for outages longer than 5 minutes (300 seconds)
locals {
  significant_outages = [
    for outage in data.uptime_outages.all.outages :
    outage if outage.duration_secs > 300
  ]
}

# Filter for outages of a specific check
locals {
  api_outages = [
    for outage in data.uptime_outages.all.outages :
    outage if can(regex("api", lower(outage.check_name)))
  ]
}

# Filter for ongoing outages (not yet resolved)
locals {
  ongoing_outages = [
    for outage in data.uptime_outages.all.outages :
    outage if !outage.state_is_up
  ]
}

# Calculate total downtime for API service
locals {
  total_api_downtime_minutes = sum([
    for outage in local.api_outages : outage.duration_secs
  ]) / 60
}

# Output outage information
output "significant_outage_count" {
  value       = length(local.significant_outages)
  description = "Number of outages longer than 5 minutes"
}

output "api_total_downtime" {
  value       = "${local.total_api_downtime_minutes} minutes"
  description = "Total downtime for API services"
}

output "ongoing_incidents" {
  value = [
    for outage in local.ongoing_outages : {
      check     = outage.check_name
      started   = outage.created_at
      locations = outage.num_locations_down
    }
  ]
  description = "Currently ongoing outages"
}
