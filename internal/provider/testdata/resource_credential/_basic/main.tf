variable display_name {
  type = string
}

variable credential_type {
  type = string
}

variable password {
  type      = string
  default   = ""
  sensitive = true
}

resource uptime_credential test {
  display_name    = var.display_name
  credential_type = var.credential_type
  secret = {
    password = var.password
  }
}
