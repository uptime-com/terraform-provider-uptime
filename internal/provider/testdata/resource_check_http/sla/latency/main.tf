variable name {
  type = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable sla_latency {
  type = string
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
  sla     = {
    latency = var.sla_latency
  }
}
