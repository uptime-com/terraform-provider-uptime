variable name {
  type = string
}

resource uptime_check_heartbeat test {
  name = var.name
}
