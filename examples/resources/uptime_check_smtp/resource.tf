# Basic SMTP check
resource "uptime_check_smtp" "example" {
  name    = "Mail Server Check"
  address = "mail.example.com"
}

# SMTP check with custom port
resource "uptime_check_smtp" "custom_port" {
  name    = "SMTP Submission Port"
  address = "mail.example.com"
  port    = 587
}

# SMTP check with encryption
resource "uptime_check_smtp" "secure" {
  name       = "Secure SMTP"
  address    = "mail.example.com"
  port       = 465
  encryption = "SSL_TLS"
}

# SMTP check with full configuration
resource "uptime_check_smtp" "full" {
  name           = "Production Mail Server"
  address        = "smtp.example.com"
  port           = 25
  interval       = 5
  contact_groups = ["nobody"]
  num_retries    = 2
}
