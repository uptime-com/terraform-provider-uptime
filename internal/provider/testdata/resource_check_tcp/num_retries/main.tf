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

variable num_retries {
  type    = number
}

resource uptime_check_tcp test {
  name        = var.name
  address     = var.address
  port        = var.port
  num_retries = var.num_retries
}
