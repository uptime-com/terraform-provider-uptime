# Basic subaccount
resource "uptime_subaccount" "example" {
  name = "Development Team"
}

# Subaccount for a specific team
resource "uptime_subaccount" "devops" {
  name        = "DevOps Team"
  description = "Monitoring resources for the DevOps team"
}

# Subaccount for a client
resource "uptime_subaccount" "client" {
  name        = "Client ABC"
  description = "Dedicated monitoring for Client ABC"
}

# Multiple subaccounts for different environments
resource "uptime_subaccount" "production" {
  name        = "Production"
  description = "Production environment monitoring"
}

resource "uptime_subaccount" "staging" {
  name        = "Staging"
  description = "Staging environment monitoring"
}

resource "uptime_subaccount" "development" {
  name        = "Development"
  description = "Development environment monitoring"
}
