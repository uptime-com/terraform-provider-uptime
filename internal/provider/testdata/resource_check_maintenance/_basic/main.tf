variable name {
  type    = string
}

variable address {
  type    = string
  default = "https://example.com"
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
}

resource uptime_check_maintenance test {
  depends_on = [uptime_check_http.test]
  check_id   = uptime_check_http.test.id
  state     = "SUPPRESSED"
}
