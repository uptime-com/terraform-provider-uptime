# Basic HTTP check
resource "uptime_check_http" "example" {
  name    = "My Website"
  address = "https://example.com"
}

# HTTP check with custom interval and contact groups
resource "uptime_check_http" "api_endpoint" {
  name           = "API Health Check"
  address        = "https://api.example.com/health"
  interval       = 5
  contact_groups = ["nobody"]
}

# HTTP check with authentication and headers
resource "uptime_check_http" "authenticated" {
  name     = "Protected Endpoint"
  address  = "https://api.example.com/admin"
  username = "monitor"
  password = "secure-password"
  headers = {
    "Authorization" = ["Bearer token123"]
    "Accept"        = ["application/json"]
  }
}

# HTTP check with content verification
resource "uptime_check_http" "content_check" {
  name               = "Homepage Content Check"
  address            = "https://www.example.com"
  expect_string      = "Welcome"
  expect_string_type = "STRING"
  status_code        = "200"
}

# HTTP check with SLA configuration
resource "uptime_check_http" "with_sla" {
  name    = "Critical Service"
  address = "https://critical.example.com"
  sla = {
    uptime  = "0.9995"
    latency = "500ms"
  }
}

# HTTP check with advanced options
resource "uptime_check_http" "advanced" {
  name                      = "Full Featured Check"
  address                   = "https://app.example.com"
  interval                  = 1
  contact_groups            = ["nobody"]
  num_retries               = 2
  sensitivity               = 2
  threshold                 = 30
  include_in_global_metrics = true
  notes                     = "Primary application endpoint"
}
