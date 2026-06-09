variable "name" {
  type = string
}

resource "uptime_contact" "test" {
  name       = var.name
  email_list = ["placeholder@example.com"]
}

resource "uptime_integration_slack" "test" {
  name           = var.name
  webhook_url    = "https://hooks.slack.com/services/test-drift"
  contact_groups = [uptime_contact.test.name]
}
