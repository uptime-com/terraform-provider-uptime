variable "name" {
  type = string
}

resource "uptime_check_http" "test" {
  name    = var.name
  address = "https://example.com"
}

resource "uptime_check_escalations" "test" {
  check_id = uptime_check_http.test.id

  escalations = [
    {
      wait_time      = 300
      num_repeats    = 3
      contact_groups = ["Default"]
    },
    {
      wait_time      = 600
      num_repeats    = 0
      contact_groups = ["Default"]
    }
  ]
}
