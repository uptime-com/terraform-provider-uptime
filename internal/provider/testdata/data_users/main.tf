variable first_name {
  type = string
}

variable last_name {
  type = string
}

variable email {
  type = string
}

variable password {
  type      = string
  sensitive = true
}

resource uptime_user test {
  first_name = var.first_name
  last_name  = var.last_name
  email      = var.email
  password   = var.password
}

data uptime_users test {
  depends_on = [uptime_user.test]
}

locals {
  filtered_users = [
    for user in data.uptime_users.test.users :
    user if user.email == var.email
  ]
}

output filtered_count {
  value = length(local.filtered_users)
}

output filtered_user_email {
  value = length(local.filtered_users) > 0 ? local.filtered_users[0].email : ""
}
