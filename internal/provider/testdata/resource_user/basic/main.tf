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

resource uptime_user test {
  email      = var.email
  first_name = var.first_name
  last_name  = var.last_name
  password   = var.password
}
