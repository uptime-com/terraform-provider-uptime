variable name {
  type = string
}

variable api_key {
  type = string
}

variable app_key {
  type = string
}

resource uptime_integration_datadog test {
  name    = var.name
  api_key = var.api_key
  app_key = var.app_key
}
