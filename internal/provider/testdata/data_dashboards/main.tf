variable name {
  type = string
}

resource uptime_dashboard test {
  name = var.name
  selected = {
    services = []
    tags     = []
  }
  services = {
    show_section = true
    sort = {
      primary   = "is_paused,cached_state_is_up"
      secondary = "-cached_last_down_alert_at"
    }
    show = {
      uptime        = true
      response_time = true
    }
  }
  alerts = {
    show_section   = false
    for_all_checks = false
  }
}

data uptime_dashboards test {
  depends_on = [uptime_dashboard.test]
}

locals {
  filtered_dashboards = [
    for dashboard in data.uptime_dashboards.test.dashboards :
    dashboard if dashboard.name == var.name
  ]
}

output filtered_count {
  value = length(local.filtered_dashboards)
}

output filtered_dashboard_name {
  value = length(local.filtered_dashboards) > 0 ? local.filtered_dashboards[0].name : ""
}
