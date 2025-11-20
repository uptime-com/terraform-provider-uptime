variable statuspage_name {
  type = string
}

resource uptime_statuspage test {
  name = var.statuspage_name
}

data uptime_statuspage_current_status test {
  statuspage_id = uptime_statuspage.test.id
}

output global_is_operational {
  value = data.uptime_statuspage_current_status.test.global_is_operational
}

output components_count {
  value = length(data.uptime_statuspage_current_status.test.components)
}
