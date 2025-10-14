# Reference an existing status page
resource "uptime_statuspage" "main" {
  name = "Production Status"
}

# Retrieve all incidents for the status page
data "uptime_statuspage_incidents" "all" {
  statuspage_id = uptime_statuspage.main.id
}

# Filter for active incidents
locals {
  active_incidents = [
    for inc in data.uptime_statuspage_incidents.all.incidents :
    inc if inc.ends_at == ""
  ]
}

# Output count of active incidents
output "active_incident_count" {
  value = length(local.active_incidents)
}
