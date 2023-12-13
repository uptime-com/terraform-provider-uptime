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

variable locations {
  type = list(string)
}

resource uptime_check_tcp test {
  name      = var.name
  address   = var.address
  port      = var.port
  locations = var.locations
}
