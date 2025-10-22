variable "name" {
  type = string
}

resource "uptime_check_http" "test" {
  name    = var.name
  address = "https://example.com"
}

resource "uptime_check_escalations" "test" {
  check_id    = uptime_check_http.test.id
  escalations = []
}
