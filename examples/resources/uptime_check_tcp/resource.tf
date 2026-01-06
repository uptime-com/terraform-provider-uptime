# Basic TCP port check
resource "uptime_check_tcp" "example" {
  name    = "Database Port Check"
  address = "db.example.com"
  port    = 5432
}

# TCP check for web server
resource "uptime_check_tcp" "web_server" {
  name    = "Web Server Port"
  address = "www.example.com"
  port    = 443
}

# TCP check with send/expect string
resource "uptime_check_tcp" "smtp_banner" {
  name          = "SMTP Banner Check"
  address       = "mail.example.com"
  port          = 25
  expect_string = "220"
}

# TCP check with full configuration
resource "uptime_check_tcp" "full" {
  name           = "Redis Port Check"
  address        = "redis.example.com"
  port           = 6379
  interval       = 5
  contact_groups = ["nobody"]
  num_retries    = 2
}
