# Create a status page first
resource "uptime_statuspage" "example" {
  name = "My Service Status"
}

# Create an incident
resource "uptime_statuspage_incident" "outage" {
  statuspage_id = uptime_statuspage.example.id
  name          = "API Degradation"
  incident_type = "INCIDENT"
  starts_at     = "2025-01-15T10:00:00Z"

  updates = [
    {
      incident_state = "investigating"
      updated_at     = "2025-01-15T10:00:00Z"
      description    = "We are investigating reports of API slowness."
    }
  ]
}

# Create a scheduled maintenance window
resource "uptime_statuspage_incident" "maintenance" {
  statuspage_id = uptime_statuspage.example.id
  name          = "Scheduled Database Maintenance"
  incident_type = "SCHEDULED_MAINTENANCE"
  starts_at     = "2025-02-01T02:00:00Z"
  ends_at       = "2025-02-01T04:00:00Z"

  updates = [
    {
      incident_state = "maintenance"
      updated_at     = "2025-02-01T02:00:00Z"
      description    = "Performing scheduled database maintenance."
    }
  ]
}
