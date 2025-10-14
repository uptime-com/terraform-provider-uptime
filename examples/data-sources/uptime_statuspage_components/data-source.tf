# Get a status page
data "uptime_statuspages" "all" {}

locals {
  prod_statuspage_id = [
    for sp in data.uptime_statuspages.all.statuspages :
    sp.id if sp.name == "Production Status"
  ][0]
}

# Retrieve all components for the status page
data "uptime_statuspage_components" "prod" {
  statuspage_id = local.prod_statuspage_id
}

# Filter for a specific component
locals {
  api_component = [
    for comp in data.uptime_statuspage_components.prod.components :
    comp if comp.name == "API Server"
  ][0]
}

# Output the component status
output "api_status" {
  value = local.api_component.status
}
