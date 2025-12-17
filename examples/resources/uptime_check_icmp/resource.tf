# Basic ICMP ping check
resource "uptime_check_icmp" "example" {
  name    = "Server Ping Check"
  address = "server.example.com"
}

# ICMP check for network device
resource "uptime_check_icmp" "router" {
  name    = "Router Ping"
  address = "192.168.1.1"
}

# ICMP check with full configuration
resource "uptime_check_icmp" "full" {
  name           = "Production Server"
  address        = "prod.example.com"
  interval       = 5
  contact_groups = ["nobody"]
  num_retries    = 3
  sensitivity    = 2
}

# ICMP check for critical infrastructure
resource "uptime_check_icmp" "critical" {
  name           = "Critical Network Node"
  address        = "core-switch.example.com"
  interval       = 1
  contact_groups = ["nobody"]
  num_retries    = 2
  sensitivity    = 1
  notes          = "Core network infrastructure"
}
