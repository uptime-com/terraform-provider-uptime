resource "uptime_integration_jira_servicedesk" "example" {
  name                         = "My JIRA Service Desk Integration"
  api_email                    = "user@example.com"
  api_token                    = "your-jira-api-token"
  jira_subdomain               = "mycompany"
  project_key                  = "SUPPORT"
  labels                       = "monitoring,uptime"
  custom_field_id_account_name = 10001
  custom_field_id_check_name   = 10002
  custom_field_id_check_url    = 10003
}
