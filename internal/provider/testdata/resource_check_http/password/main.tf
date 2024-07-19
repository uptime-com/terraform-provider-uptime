variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable password {
  type = string
  sensitive = false
}

resource uptime_check_http test {
  name     = var.name
  address  = var.address
  password = var.password
}
