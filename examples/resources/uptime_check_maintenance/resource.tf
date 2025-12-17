# Basic weekly maintenance window
resource "uptime_check_http" "api" {
  name    = "API Health"
  address = "https://api.example.com/health"
}

resource "uptime_check_maintenance" "api" {
  check_id = uptime_check_http.api.id
  state    = "SCHEDULED"
  schedule = [
    {
      type      = "WEEKLY"
      weekdays  = [0] # Sunday
      from_time = "02:00:00"
      to_time   = "04:00:00"
    }
  ]
}

# Suppress alerts for a check (no schedule needed)
resource "uptime_check_http" "website" {
  name    = "Website Check"
  address = "https://www.example.com"
}

resource "uptime_check_maintenance" "website" {
  check_id = uptime_check_http.website.id
  state    = "SUPPRESSED"
}

# Weekly maintenance with pause option
resource "uptime_check_http" "app" {
  name    = "App Health"
  address = "https://app.example.com/health"
}

resource "uptime_check_maintenance" "app" {
  check_id                       = uptime_check_http.app.id
  state                          = "SCHEDULED"
  pause_on_scheduled_maintenance = true
  schedule = [
    {
      type      = "WEEKLY"
      weekdays  = [0, 6] # Saturday and Sunday
      from_time = "02:00:00"
      to_time   = "06:00:00"
    }
  ]
}
