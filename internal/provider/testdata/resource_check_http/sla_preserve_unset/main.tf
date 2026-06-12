variable name { type = string }

resource uptime_check_http test {
  name    = var.name
  address = "https://example.com"
}
