# Create a status page first
resource "uptime_statuspage" "example" {
  name             = "My Service Status"
  visibility_level = "EXTERNAL_USERS"
}

# Add a user who can access the status page
resource "uptime_statuspage_user" "admin" {
  statuspage_id = uptime_statuspage.example.id
  email         = "admin@example.com"
  first_name    = "John"
  last_name     = "Admin"
  is_active     = true
}
