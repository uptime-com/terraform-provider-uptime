variable statuspage_name {
  type = string
}

variable incident_name {
  type = string
}

resource uptime_statuspage test {
  name = var.statuspage_name
}

resource uptime_statuspage_incident test {
  statuspage_id = uptime_statuspage.test.id
  name          = var.incident_name
  starts_at     = "2024-01-01T00:00:00Z"
  incident_type = "INCIDENT"
  updates = [{
    message        = "Initial incident update"
    incident_state = "investigating"
    notify         = false
    updated_at     = "2024-01-01T00:00:00Z"
  }]
}

data uptime_statuspage_incidents test {
  statuspage_id = uptime_statuspage.test.id
  depends_on    = [uptime_statuspage_incident.test]
}

locals {
  filtered_incidents = [
    for inc in data.uptime_statuspage_incidents.test.incidents :
    inc if inc.name == var.incident_name
  ]
}

output filtered_count {
  value = length(local.filtered_incidents)
}

output filtered_incident_name {
  value = length(local.filtered_incidents) > 0 ? local.filtered_incidents[0].name : ""
}
