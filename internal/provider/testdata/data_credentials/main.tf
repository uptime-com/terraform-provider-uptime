variable display_name {
  type = string
}

variable password {
  type      = string
  sensitive = true
}

resource uptime_credential test {
  display_name    = var.display_name
  credential_type = "BASIC"
  secret = {
    password = var.password
  }
}

data uptime_credentials test {
  depends_on = [uptime_credential.test]
}

locals {
  filtered_credentials = [
    for cred in data.uptime_credentials.test.credentials :
    cred if cred.display_name == var.display_name
  ]
}

output filtered_count {
  value = length(local.filtered_credentials)
}

output filtered_credential_name {
  value = length(local.filtered_credentials) > 0 ? local.filtered_credentials[0].display_name : ""
}

output filtered_credential_type {
  value = length(local.filtered_credentials) > 0 ? local.filtered_credentials[0].credential_type : ""
}
