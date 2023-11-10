variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable config_crl {
  type    = bool
}

resource uptime_check_sslcert test {
  name    = var.name
  address = var.address
  config  = {
    crl = var.config_crl
  }
}
