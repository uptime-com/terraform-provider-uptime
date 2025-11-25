variable name {
  type = string
}

variable webhook_url {
  type = string
}

resource uptime_integration_microsoft_teams test {
  name        = var.name
  webhook_url = var.webhook_url
}
