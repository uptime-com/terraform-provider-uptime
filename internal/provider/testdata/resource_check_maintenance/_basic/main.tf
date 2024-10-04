variable name {
  type    = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable state {
  type    = string
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
}

resource uptime_check_maintenance test {
  check_id   = uptime_check_http.test.id
  state      = var.state
}
