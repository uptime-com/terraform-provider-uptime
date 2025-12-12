# Create a basic contact
resource "uptime_contact" "example" {
  name = "Operations Team"
}

# Create a contact with email and SMS notifications
resource "uptime_contact" "ops" {
  name       = "Operations Team"
  email_list = ["ops@example.com", "alerts@example.com"]
  sms_list   = ["+1234567890"]
}

# Create a contact with integrations
resource "uptime_contact" "with_integrations" {
  name         = "DevOps Team"
  email_list   = ["devops@example.com"]
  integrations = ["slack-alerts", "pagerduty-oncall"]
}
