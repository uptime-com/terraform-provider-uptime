# Create a status page first
resource "uptime_statuspage" "example" {
  name = "My Service Status"
}

# Create an HTTP check to link to the component
resource "uptime_check_http" "api" {
  name    = "API Health"
  address = "https://api.example.com/health"
}

# Add a component to the status page
resource "uptime_statuspage_component" "api" {
  statuspage_id = uptime_statuspage.example.id
  name          = "API Service"
  description   = "Core API endpoints"
  service_id    = uptime_check_http.api.id
}
