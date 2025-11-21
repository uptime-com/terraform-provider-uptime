variable name {
  type = string
}

variable email {
  type = string
}

resource uptime_integration_pushbullet test {
  name  = var.name
  email = var.email
}
