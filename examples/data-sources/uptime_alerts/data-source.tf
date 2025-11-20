# Retrieve alerts from your Uptime.com account
data "uptime_alerts" "all" {}

# Filter for unresolved (down) alerts
locals {
  active_alerts = [
    for alert in data.uptime_alerts.all.alerts :
    alert if !alert.state_is_up && !alert.ignored
  ]
}

# Filter alerts for a specific check
locals {
  api_alerts = [
    for alert in data.uptime_alerts.all.alerts :
    alert if can(regex("api", lower(alert.check_name)))
  ]
}

# Filter for recent alerts (those with resolved_at empty means still active)
locals {
  unresolved_alerts = [
    for alert in data.uptime_alerts.all.alerts :
    alert if alert.resolved_at == ""
  ]
}

# Output alert information
output "active_alert_count" {
  value       = length(local.active_alerts)
  description = "Number of active (down, not ignored) alerts"
}

output "api_service_alerts" {
  value = [
    for alert in local.api_alerts : {
      check   = alert.check_name
      output  = alert.output
      created = alert.created_at
    }
  ]
  description = "Recent alerts for API services"
}
