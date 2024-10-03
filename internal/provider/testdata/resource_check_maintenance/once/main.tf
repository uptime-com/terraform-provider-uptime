variable name {
  type    = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable once_start_date {
  type    = string
}

variable once_end_date {
  type    = string
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
}

resource uptime_check_maintenance test {
  check_id   = uptime_check_http.test.id
  state      = "SCHEDULED"
  schedule   = [{
    type = "ONCE"
    once_start_date = var.once_start_date
    once_end_date   = var.once_end_date
  }]
}
