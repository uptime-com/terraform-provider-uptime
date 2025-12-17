# Basic IMAP check
resource "uptime_check_imap" "example" {
  name    = "IMAP Server Check"
  address = "mail.example.com"
}

# IMAP check with custom port
resource "uptime_check_imap" "custom_port" {
  name    = "IMAP Custom Port"
  address = "mail.example.com"
  port    = 993
}

# IMAP check with encryption
resource "uptime_check_imap" "secure" {
  name       = "Secure IMAP"
  address    = "mail.example.com"
  port       = 993
  encryption = "SSL_TLS"
}

# IMAP check with full configuration
resource "uptime_check_imap" "full" {
  name           = "Production IMAP"
  address        = "imap.example.com"
  port           = 993
  encryption     = "SSL_TLS"
  interval       = 5
  contact_groups = ["nobody"]
  num_retries    = 2
}
