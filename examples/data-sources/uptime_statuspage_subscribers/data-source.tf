# Reference an existing status page
resource "uptime_statuspage" "main" {
  name = "Production Status"
}

# Retrieve all subscribers for the status page
data "uptime_statuspage_subscribers" "all" {
  statuspage_id = uptime_statuspage.main.id
}

# Filter subscribers by type
locals {
  email_subscribers = [
    for sub in data.uptime_statuspage_subscribers.all.subscribers :
    sub if sub.type == "email"
  ]

  sms_subscribers = [
    for sub in data.uptime_statuspage_subscribers.all.subscribers :
    sub if sub.type == "sms"
  ]
}

# Output subscriber counts
output "email_subscriber_count" {
  value = length(local.email_subscribers)
}

output "sms_subscriber_count" {
  value = length(local.sms_subscribers)
}
