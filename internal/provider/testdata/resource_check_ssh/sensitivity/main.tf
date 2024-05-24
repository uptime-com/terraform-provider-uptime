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

variable sensitivity {
  type    = string
}

resource uptime_check_ssh test {
  name          = var.name
  address       = var.address
  port          = var.port
  sensitivity   = var.sensitivity
}
