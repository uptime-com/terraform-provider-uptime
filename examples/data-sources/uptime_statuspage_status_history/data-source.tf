# Use with a status page resource
resource "uptime_statuspage" "main" {
  name = "Service Status"
}

# Retrieve status history for a status page
data "uptime_statuspage_status_history" "main" {
  statuspage_id = uptime_statuspage.main.id
}

# Output status history
output "status_history" {
  value       = data.uptime_statuspage_status_history.main
  description = "Historical status data for the status page"
}

output "history_entries" {
  value       = data.uptime_statuspage_status_history.main.history
  description = "List of historical status entries"
}
