variable name {
  type = string
}

variable allow_subscriptions_email {
  type = bool
}

variable allow_subscriptions_rss {
  type = bool
}

variable allow_subscriptions_slack {
  type = bool
}

variable allow_subscriptions_sms {
  type = bool
}

variable allow_subscriptions_webhook {
  type = bool
}

variable hide_empty_tabs_history {
  type = bool
}

variable theme {
  type = string
}

variable custom_header_bg_color_hex {
  type = string
}

variable custom_header_text_color_hex {
  type = string
}

resource "uptime_statuspage" "test" {
  name                       = var.name
  allow_subscriptions_email  = var.allow_subscriptions_email
  allow_subscriptions_rss    = var.allow_subscriptions_rss
  allow_subscriptions_slack  = var.allow_subscriptions_slack
  allow_subscriptions_sms    = var.allow_subscriptions_sms
  allow_subscriptions_webhook = var.allow_subscriptions_webhook
  hide_empty_tabs_history    = var.hide_empty_tabs_history
  theme                      = var.theme
  custom_header_bg_color_hex = var.custom_header_bg_color_hex
  custom_header_text_color_hex = var.custom_header_text_color_hex
}
