# Basic DNS check
resource "uptime_check_dns" "example" {
  name    = "DNS Resolution Check"
  address = "example.com"
}

# DNS check with record type validation
resource "uptime_check_dns" "mx_record" {
  name            = "MX Record Check"
  address         = "example.com"
  dns_record_type = "MX"
}

# DNS check for specific nameserver
resource "uptime_check_dns" "custom_ns" {
  name            = "Custom NS Check"
  address         = "example.com"
  dns_server      = "8.8.8.8"
  dns_record_type = "A"
}

# DNS check with expected IP validation
resource "uptime_check_dns" "validate_ip" {
  name            = "DNS IP Validation"
  address         = "www.example.com"
  dns_record_type = "A"
  expect_string   = "93.184.215.14"
}

# DNS check with full configuration
resource "uptime_check_dns" "full" {
  name            = "Production DNS"
  address         = "api.example.com"
  dns_record_type = "A"
  dns_server      = "8.8.8.8"
  interval        = 5
  contact_groups  = ["nobody"]
  num_retries     = 2
}
