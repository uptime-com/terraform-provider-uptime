# Basic tag
resource "uptime_tag" "example" {
  tag       = "production"
  color_hex = "#2ecc71"
}

# Tag with custom color
resource "uptime_tag" "environment" {
  tag       = "staging"
  color_hex = "#ff5733"
}

# Tags for different purposes
resource "uptime_tag" "team" {
  tag       = "devops-team"
  color_hex = "#3498db"
}

resource "uptime_tag" "critical" {
  tag       = "critical"
  color_hex = "#e74c3c"
}

resource "uptime_tag" "customer_facing" {
  tag       = "customer-facing"
  color_hex = "#27ae60"
}
