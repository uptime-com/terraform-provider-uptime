# Basic webhook check (incoming webhook)
resource "uptime_check_webhook" "example" {
  name = "Application Webhook"
}

# Webhook check with full configuration
resource "uptime_check_webhook" "full" {
  name           = "Critical App Webhook"
  contact_groups = ["nobody"]
  sla = {
    uptime = "0.999"
  }
}

# Output the webhook URL to use in your application
output "webhook_url" {
  value       = uptime_check_webhook.example.url
  description = "URL for your application to POST status updates"
}
