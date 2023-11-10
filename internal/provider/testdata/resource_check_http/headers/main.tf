variable name {
  type = string
}

variable address {
  type    = string
  default = "https://example.com"
}

variable headers {
  type = map(list(string))
}

resource uptime_check_http test {
  name    = var.name
  address = var.address
  headers = var.headers
}
