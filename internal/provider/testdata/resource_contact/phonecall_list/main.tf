resource uptime_contact test {
  name           = var.name
  phonecall_list = var.phonecall_list
}

variable name {
  type = string
}

variable phonecall_list {
  type = list(string)
}
