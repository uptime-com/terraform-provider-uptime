variable name {
  type = string
}

variable service_name {
  type    = string
  default = "aws-ec2-us-east-1"
}

resource uptime_check_cloudstatus test {
  name         = var.name
  service_name = var.service_name
}
