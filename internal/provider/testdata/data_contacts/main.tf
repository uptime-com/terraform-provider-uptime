variable name {
  type = string
}

variable email {
  type = string
}

resource uptime_contact test {
  name       = var.name
  email_list = [var.email]
}

data uptime_contacts test {
  depends_on = [uptime_contact.test]
}

locals {
  filtered_contacts = [
    for contact in data.uptime_contacts.test.contacts :
    contact if contact.name == var.name
  ]
}

output filtered_count {
  value = length(local.filtered_contacts)
}

output filtered_contact_name {
  value = length(local.filtered_contacts) > 0 ? local.filtered_contacts[0].name : ""
}
