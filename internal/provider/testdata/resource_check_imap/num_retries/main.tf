variable name {
  type = string
}

variable address {
  type    = string
  default = "example.com"
}

variable num_retries {
  type = number
}

resource uptime_check_imap test {
  name        = var.name
  address     = var.address
  num_retries = var.num_retries
}
