# Create a status page first
resource "uptime_statuspage" "example" {
  name                      = "My Service Status"
  allow_subscriptions       = true
  allow_subscriptions_email = true
  allow_subscriptions_webhook = true
}

# Add an email subscriber
resource "uptime_statuspage_subscriber" "email" {
  statuspage_id = uptime_statuspage.example.id
  type          = "EMAIL"
  target        = "alerts@example.com"
}

# Add a webhook subscriber
resource "uptime_statuspage_subscriber" "webhook" {
  statuspage_id = uptime_statuspage.example.id
  type          = "WEBHOOK"
  target        = "https://example.com/status-webhook"
}
