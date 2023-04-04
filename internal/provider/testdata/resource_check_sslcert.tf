resource uptime_check_sslcert test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody"]
  address        = "example.com"
  config {
    crl = false
  }
}
// ---
resource uptime_check_sslcert test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody", "noone"]
  address        = "example.net"
  config {
    crl = true
  }
}
