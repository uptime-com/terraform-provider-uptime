variable name {
  type = string
}

variable email {
  type = string
}

variable api_token {
  type = string
}

variable metric_name {
  type = string
}

resource uptime_integration_librato test {
  name        = var.name
  email       = var.email
  api_token   = var.api_token
  metric_name = var.metric_name
}
