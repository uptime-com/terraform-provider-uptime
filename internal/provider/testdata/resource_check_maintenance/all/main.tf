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

variable monthly_from_time {
  type    = string
}

variable monthly_to_time {
  type    = string
}

variable monthday {
  type    = number
  default = 0
}

variable monthday_from {
  type    = number
  default = 0
}

variable monthday_to {
  type    = number
  default = 0
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
  schedule   = [
    {
      type = "WEEKLY"
      from_time = var.from_time
      to_time   = var.to_time
      weekdays  = var.weekdays
    },
    {
      type = "MONTHLY"
      from_time     = var.monthly_from_time
      to_time       = var.monthly_to_time
      monthday      = var.monthday
      monthday_from = var.monthday_from
      monthday_to   = var.monthday_to
    },
    {
      type = "ONCE"
      once_start_date = var.once_start_date
      once_end_date   = var.once_end_date
    }
  ]
}
