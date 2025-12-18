# Use with a status page resource
resource "uptime_statuspage" "main" {
  name = "Service Status"
}

# Retrieve current status of a status page
data "uptime_statuspage_current_status" "main" {
  statuspage_id = uptime_statuspage.main.id
}

# Output current status information
output "current_status" {
  value       = data.uptime_statuspage_current_status.main
  description = "Current status of the status page"
}

output "is_operational" {
  value       = data.uptime_statuspage_current_status.main.global_is_operational
  description = "Whether all services are operational"
}
