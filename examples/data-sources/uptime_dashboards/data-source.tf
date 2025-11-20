# Retrieve all dashboards configured in your Uptime.com account
data "uptime_dashboards" "all" {}

# Filter for a specific dashboard by name
locals {
  main_dashboard = [
    for dashboard in data.uptime_dashboards.all.dashboards :
    dashboard if dashboard.name == "Production Overview"
  ][0]
}

# Filter for pinned dashboards
locals {
  pinned_dashboards = [
    for dashboard in data.uptime_dashboards.all.dashboards :
    dashboard if dashboard.is_pinned
  ]
}

# Output dashboard information
output "main_dashboard_id" {
  value       = local.main_dashboard.id
  description = "ID of the main production dashboard"
}

output "pinned_dashboard_names" {
  value       = [for d in local.pinned_dashboards : d.name]
  description = "Names of all pinned dashboards"
}
