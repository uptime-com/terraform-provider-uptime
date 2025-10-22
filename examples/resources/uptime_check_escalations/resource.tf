# Create an HTTP check
resource "uptime_check_http" "website" {
  name    = "Website Health Check"
  address = "https://example.com"
}

# Configure escalation rules for the check
resource "uptime_check_escalations" "website" {
  check_id = uptime_check_http.website.id

  escalations = [
    {
      # First escalation: wait 5 minutes, notify Default group 3 times
      wait_time      = 300
      num_repeats    = 3
      contact_groups = ["Default"]
    },
    {
      # Second escalation: wait 10 more minutes, notify On-Call team indefinitely
      wait_time      = 600
      num_repeats    = 0 # 0 means repeat indefinitely
      contact_groups = ["On-Call Team"]
    }
  ]
}
