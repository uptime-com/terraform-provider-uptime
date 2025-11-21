variable name {
  type = string
}

variable user {
  type = string
}

resource uptime_integration_pushover test {
  name = var.name
  user = var.user
}
