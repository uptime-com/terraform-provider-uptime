# Basic dashboard
resource "uptime_dashboard" "example" {
  name   = "My Dashboard"
  alerts = {}
  services = {
    show = {}
    sort = {}
  }
  selected = {}
}

# Dashboard with service filtering options
resource "uptime_dashboard" "configured" {
  name = "Production Services"
  alerts = {
    show_section = true
  }
  services = {
    show_section = true
    include = {
      up          = true
      down        = true
      paused      = false
      maintenance = true
    }
    show = {
      uptime        = true
      response_time = true
    }
    sort = {
      primary   = "is_paused,cached_state_is_up"
      secondary = "-cached_last_down_alert_at"
    }
  }
  selected = {}
}

# Dashboard using check IDs from created resources
resource "uptime_check_http" "api" {
  name    = "API Health"
  address = "https://api.example.com/health"
}

resource "uptime_check_http" "web" {
  name    = "Web Frontend"
  address = "https://www.example.com"
}

resource "uptime_dashboard" "dynamic" {
  name = "Dynamic Dashboard"
  alerts = {}
  services = {
    show = {}
    sort = {}
  }
  selected = {
    services = [uptime_check_http.api.name, uptime_check_http.web.name]
  }
}
