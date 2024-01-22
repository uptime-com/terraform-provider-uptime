variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable port {
  type = number
}

resource uptime_check_imap test {
  name    = var.name
  address = var.address
  port    = var.port
}
