variable statuspage_name {
  type = string
}

variable email {
  type = string
}

resource uptime_statuspage test {
  name = var.statuspage_name
}

resource uptime_statuspage_subscriber test {
  statuspage_id = uptime_statuspage.test.id
  target        = var.email
  type          = "EMAIL"
}

data uptime_statuspage_subscribers test {
  statuspage_id = uptime_statuspage.test.id
  depends_on    = [uptime_statuspage_subscriber.test]
}

locals {
  filtered_subscribers = [
    for sub in data.uptime_statuspage_subscribers.test.subscribers :
    sub if sub.target == var.email
  ]
}

output filtered_count {
  value = length(local.filtered_subscribers)
}

output filtered_subscriber_target {
  value = length(local.filtered_subscribers) > 0 ? local.filtered_subscribers[0].target : ""
}
