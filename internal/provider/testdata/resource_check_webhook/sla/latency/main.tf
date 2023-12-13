variable name {
  type = string
}

variable sla_latency {
  type = string
}

resource uptime_check_webhook test {
  name   = var.name
  sla    = {
    latency = var.sla_latency
  }
}
