variable name {
  type    = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable from_time {
  type    = string
}

variable to_time {
  type    = string
}

variable weekdays {
  type    = list(number)
  default = []
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
}

resource uptime_check_maintenance test {
  check_id   = uptime_check_http.test.id
  state      = "SCHEDULED"
  schedule   = [{
    type = "WEEKLY"
    from_time = var.from_time
    to_time   = var.to_time
    weekdays  = var.weekdays
  }]
}
