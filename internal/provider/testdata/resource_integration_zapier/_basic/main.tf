variable name {
  type = string
}

variable webhook_url {
  type = string
}

resource uptime_integration_zapier test {
  name        = var.name
  webhook_url = var.webhook_url
}
