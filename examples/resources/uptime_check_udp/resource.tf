# Basic UDP port check
resource "uptime_check_udp" "example" {
  name          = "DNS UDP Check"
  address       = "dns.example.com"
  port          = 53
  send_string   = "query"
  expect_string = "response"
}

# UDP check with send/expect
resource "uptime_check_udp" "with_payload" {
  name          = "UDP Service Check"
  address       = "service.example.com"
  port          = 5000
  send_string   = "ping"
  expect_string = "pong"
}

# UDP check with full configuration
resource "uptime_check_udp" "full" {
  name           = "Game Server Check"
  address        = "game.example.com"
  port           = 27015
  send_string    = "status"
  expect_string  = "ok"
  interval       = 5
  contact_groups = ["nobody"]
  num_retries    = 2
}
