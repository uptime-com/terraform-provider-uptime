variable name {
  type = string
}

variable user {
  type = string
}

variable priority {
  type = number
}

resource uptime_integration_pushover test {
  name     = var.name
  user     = var.user
  priority = var.priority
}
