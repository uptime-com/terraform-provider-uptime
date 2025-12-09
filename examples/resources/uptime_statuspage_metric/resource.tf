# Create a status page first
resource "uptime_statuspage" "example" {
  name = "My Service Status"
}

# Create an HTTP check to use as metric source
resource "uptime_check_http" "api" {
  name    = "API Health"
  address = "https://api.example.com/health"
}

# Add a metric to display on the status page
resource "uptime_statuspage_metric" "api_response_time" {
  statuspage_id = uptime_statuspage.example.id
  name          = "API Response Time"
  service_id    = uptime_check_http.api.id
  is_visible    = true
}
