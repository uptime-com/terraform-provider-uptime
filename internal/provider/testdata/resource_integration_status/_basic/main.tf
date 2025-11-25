variable name {
  type = string
}

variable statuspage_id {
  type = string
}

variable api_id {
  type = string
}

variable api_key {
  type = string
}

variable metric {
  type = string
}

resource uptime_integration_status test {
  name          = var.name
  statuspage_id = var.statuspage_id
  api_id        = var.api_id
  api_key       = var.api_key
  metric        = var.metric
}
