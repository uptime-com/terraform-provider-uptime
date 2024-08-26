variable "name" {
  type = string
}

variable "address" {
  type    = string
  default = "example.com"
}

variable "port" {
  type = number
}

variable "encryption" {
  type    = string
  default = "" # or "SSL_TLS"
}

resource "uptime_check_tcp" "test" {
  name       = var.name
  address    = var.address
  port       = var.port
  encryption = var.encryption
}
