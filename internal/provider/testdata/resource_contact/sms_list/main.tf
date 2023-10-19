resource uptime_contact test {
  name           = var.name
  sms_list       = var.sms_list
}

variable name {
  type = string
}

variable sms_list {
  type = list(string)
}
