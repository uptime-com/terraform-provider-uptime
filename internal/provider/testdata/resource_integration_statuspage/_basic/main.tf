variable name {
  type = string
}

variable api_key {
  type = string
}

variable page {
  type = string
}

variable metric {
  type = string
}

resource uptime_integration_statuspage test {
  name    = var.name
  api_key = var.api_key
  page    = var.page
  metric  = var.metric
}
