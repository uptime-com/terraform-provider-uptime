# Basic SSH check
resource "uptime_check_ssh" "example" {
  name    = "SSH Server Check"
  address = "server.example.com"
  port    = 22
}

# SSH check with custom port
resource "uptime_check_ssh" "custom_port" {
  name    = "SSH Custom Port"
  address = "server.example.com"
  port    = 2222
}

# SSH check with full configuration
resource "uptime_check_ssh" "full" {
  name           = "Production SSH"
  address        = "prod.example.com"
  port           = 22
  interval       = 5
  contact_groups = ["nobody"]
  num_retries    = 2
}
