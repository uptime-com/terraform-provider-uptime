variable name {
  type = string
}

variable email_list {
  type = list(string)
}

resource uptime_contact test {
  name           = var.name
  email_list     = var.email_list
}
