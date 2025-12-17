# Basic NTP check
resource "uptime_check_ntp" "example" {
  name    = "NTP Server Check"
  address = "time.example.com"
}

# NTP check for public server
resource "uptime_check_ntp" "public" {
  name    = "Public NTP"
  address = "pool.ntp.org"
}

# NTP check with full configuration
resource "uptime_check_ntp" "full" {
  name           = "Production NTP"
  address        = "ntp.example.com"
  interval       = 30
  contact_groups = ["nobody"]
  num_retries    = 2
}
