variable email {
  type = string
}

variable first_name {
  type = string
}

variable last_name {
  type = string
}

variable password {
  type      = string
  sensitive = true
}

variable is_api_enabled {
  type = bool
}

variable notify_paid_invoices {
  type = bool
}

resource uptime_user test {
  email                = var.email
  first_name           = var.first_name
  last_name            = var.last_name
  password             = var.password
  is_api_enabled       = var.is_api_enabled
  notify_paid_invoices = var.notify_paid_invoices
}
