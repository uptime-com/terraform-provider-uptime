variable name {
  type = string
}

variable sla_uptime {
  type = string
}

resource uptime_check_heartbeat test {
  name   = var.name
  sla    = {
    uptime = var.sla_uptime
  }
}
