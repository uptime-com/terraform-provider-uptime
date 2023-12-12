variable "name" {
  type = string
}

resource "uptime_check_group" "test" {
  name   = var.name
  config = {
    response_time = {
      check_type = "HTTP"
    }
  }
}
