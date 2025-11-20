variable statuspage_name {
  type = string
}

resource uptime_statuspage test {
  name = var.statuspage_name
}

data uptime_statuspage_status_history test {
  statuspage_id = uptime_statuspage.test.id
}

output history_count {
  value = length(data.uptime_statuspage_status_history.test.history)
}
