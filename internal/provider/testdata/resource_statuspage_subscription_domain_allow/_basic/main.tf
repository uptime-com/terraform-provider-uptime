variable "name" {
  type = string
}

variable "domain" {
  type = string
}

resource "uptime_statuspage" "test" {
  name = var.name
}

resource "uptime_statuspage_subscription_domain_allow" "test" {
  statuspage_id = uptime_statuspage.test.id
  domain        = var.domain
}
