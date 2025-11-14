variable "name" {
  type = string
}

variable "address" {
  type    = string
  default = "example.com"
}

variable "resolve" {
  type    = string
  default = ""
}

resource "uptime_check_sslcert" "test" {
  name    = var.name
  address = var.address

  config = {
    resolve                    = var.resolve
    ignore_authority_warnings  = true
    ignore_sct                 = false
  }
}
