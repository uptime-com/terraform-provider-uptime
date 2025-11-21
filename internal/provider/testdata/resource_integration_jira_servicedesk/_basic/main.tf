variable name {
  type = string
}

variable api_email {
  type = string
}

variable api_token {
  type = string
}

variable jira_subdomain {
  type = string
}

variable project_key {
  type = string
}

resource uptime_integration_jira_servicedesk test {
  name           = var.name
  api_email      = var.api_email
  api_token      = var.api_token
  jira_subdomain = var.jira_subdomain
  project_key    = var.project_key
}
