variable name {
  type = string
}

variable interval {
  type = number
}

resource uptime_check_heartbeat test {
  name     = var.name
  interval = var.interval
}
