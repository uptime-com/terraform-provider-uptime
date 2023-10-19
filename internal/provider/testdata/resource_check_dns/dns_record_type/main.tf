variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable dns_record_type {
  type = string
}

resource uptime_check_dns test {
  name            = var.name
  address         = var.address
  dns_record_type = var.dns_record_type
}
