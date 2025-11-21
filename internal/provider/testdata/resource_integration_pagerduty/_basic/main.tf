variable name {
  type = string
}

variable service_key {
  type = string
}

resource uptime_integration_pagerduty test {
  name        = var.name
  service_key = var.service_key
}
