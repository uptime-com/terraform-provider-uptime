# Basic POP3 check
resource "uptime_check_pop" "example" {
  name    = "POP3 Server Check"
  address = "mail.example.com"
}

# POP3 check with custom port
resource "uptime_check_pop" "custom_port" {
  name    = "POP3 Custom Port"
  address = "mail.example.com"
  port    = 995
}

# POP3 check with encryption
resource "uptime_check_pop" "secure" {
  name       = "Secure POP3"
  address    = "mail.example.com"
  port       = 995
  encryption = "SSL_TLS"
}

# POP3 check with full configuration
resource "uptime_check_pop" "full" {
  name           = "Production POP3"
  address        = "pop.example.com"
  port           = 995
  encryption     = "SSL_TLS"
  interval       = 5
  contact_groups = ["nobody"]
  num_retries    = 2
}
