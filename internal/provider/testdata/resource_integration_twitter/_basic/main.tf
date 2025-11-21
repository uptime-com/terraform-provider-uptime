variable name {
  type = string
}

variable oauth_token {
  type = string
}

variable oauth_token_secret {
  type = string
}

resource uptime_integration_twitter test {
  name               = var.name
  oauth_token        = var.oauth_token
  oauth_token_secret = var.oauth_token_secret
}
