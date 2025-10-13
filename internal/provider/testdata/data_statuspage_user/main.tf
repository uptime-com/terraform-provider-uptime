variable statuspage_name {
  type = string
}

variable email {
  type = string
}

variable first_name {
  type = string
}

variable last_name {
  type = string
}

resource uptime_statuspage test {
  name = var.statuspage_name
}

resource uptime_statuspage_user test {
  statuspage_id = uptime_statuspage.test.id
  email         = var.email
  first_name    = var.first_name
  last_name     = var.last_name
  is_active     = true
}

data uptime_statuspage_users test {
  statuspage_id = uptime_statuspage.test.id
  depends_on    = [uptime_statuspage_user.test]
}

locals {
  filtered_users = [
    for user in data.uptime_statuspage_users.test.users :
    user if user.email == var.email
  ]
}

output filtered_count {
  value = length(local.filtered_users)
}

output filtered_user_email {
  value = length(local.filtered_users) > 0 ? local.filtered_users[0].email : ""
}
