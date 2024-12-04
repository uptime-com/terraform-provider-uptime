variable name {
  type = string
}

variable api_endpoint {
  type = string
}

variable api_key {
  type = string
}

resource uptime_integration_opsgenie test {
  name         = var.name
  api_endpoint = var.api_endpoint
  api_key      = var.api_key
}
