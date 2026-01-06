# Basic heartbeat check
resource "uptime_check_heartbeat" "example" {
  name = "Cron Job Monitor"
}

# Heartbeat check for scheduled task
resource "uptime_check_heartbeat" "backup" {
  name     = "Daily Backup Monitor"
  interval = 1440 # Expect signal every 24 hours
}

# Heartbeat check with contact groups
resource "uptime_check_heartbeat" "etl" {
  name           = "ETL Pipeline"
  interval       = 60
  contact_groups = ["nobody"]
}

# Heartbeat check with full configuration
resource "uptime_check_heartbeat" "full" {
  name           = "Critical Cron Job"
  interval       = 30
  contact_groups = ["nobody"]
  notes          = "Monitor for hourly data sync job"
  sla = {
    uptime = "0.999"
  }
}

# Output the heartbeat URL to use in your application
output "heartbeat_url" {
  value       = uptime_check_heartbeat.example.url
  description = "URL to ping from your scheduled task"
}
