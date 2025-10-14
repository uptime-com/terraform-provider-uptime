# Retrieve all status pages
data "uptime_statuspages" "all" {}

# Filter for a specific status page by name
locals {
  my_statuspage = [
    for sp in data.uptime_statuspages.all.statuspages :
    sp if sp.name == "Production Status"
  ][0]
}

# Use the status page ID to create a component
resource "uptime_statuspage_component" "api" {
  statuspage_id = local.my_statuspage.id
  name          = "API Server"
  description   = "Main API endpoint"
}
