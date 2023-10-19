variable name {
  type = string
}

variable response_time_sla {
  type = string
}

resource uptime_check_heartbeat test {
  name              = var.name
  response_time_sla = var.response_time_sla
}
