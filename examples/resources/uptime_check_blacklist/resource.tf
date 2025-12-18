# Basic blacklist check (email/IP reputation)
resource "uptime_check_blacklist" "example" {
  name    = "IP Blacklist Check"
  address = "mail.example.com"
}

# Blacklist check for mail server
resource "uptime_check_blacklist" "mail" {
  name    = "Mail Server Reputation"
  address = "smtp.example.com"
}

# Blacklist check with full configuration
resource "uptime_check_blacklist" "full" {
  name           = "Production Mail Reputation"
  address        = "mail.example.com"
  contact_groups = ["nobody"]
  num_retries    = 2
}
