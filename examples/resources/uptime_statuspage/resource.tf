# Create a basic status page
resource "uptime_statuspage" "example" {
  name = "My Service Status"
}

# Create a public status page with custom settings
resource "uptime_statuspage" "public" {
  name                      = "Public Status Page"
  page_type                 = "PUBLIC"
  visibility_level          = "PUBLIC"
  allow_subscriptions       = true
  allow_subscriptions_email = true
  allow_subscriptions_rss   = true
  show_status_tab           = true
  show_history_tab          = true
  show_active_incidents     = true
  timezone                  = "America/New_York"
}
