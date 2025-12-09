# Create a status page with email subscriptions enabled
resource "uptime_statuspage" "example" {
  name                      = "My Service Status"
  allow_subscriptions       = true
  allow_subscriptions_email = true
}

# Block subscriptions from specific domains
resource "uptime_statuspage_subscription_domain_block" "spam" {
  statuspage_id = uptime_statuspage.example.id
  domain        = "spam.com"
}

resource "uptime_statuspage_subscription_domain_block" "competitor" {
  statuspage_id = uptime_statuspage.example.id
  domain        = "competitor.com"
}
