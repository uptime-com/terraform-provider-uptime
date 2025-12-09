# Create a status page with email subscriptions enabled
resource "uptime_statuspage" "example" {
  name                      = "My Service Status"
  allow_subscriptions       = true
  allow_subscriptions_email = true
}

# Allow subscriptions only from specific domains
resource "uptime_statuspage_subscription_domain_allow" "company" {
  statuspage_id = uptime_statuspage.example.id
  domain        = "example.com"
}

resource "uptime_statuspage_subscription_domain_allow" "partner" {
  statuspage_id = uptime_statuspage.example.id
  domain        = "partner.com"
}
