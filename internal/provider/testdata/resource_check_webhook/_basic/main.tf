variable "name" {
  type = string
}

resource "uptime_check_webhook" test {
  name = var.name
}
