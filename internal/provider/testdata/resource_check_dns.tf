resource uptime_check_dns test {
  name    = "{{ petname 3 "-" }}"
  address = "example.com"
}
// ---
resource uptime_check_dns test {
  name              = "{{ petname 3 "-" }}"
  contact_groups    = ["nobody", "noone"]
  interval          = 10
  locations         = ["Serbia", "Austria"]
  address           = "example.net"
  dns_record_type   = "AAAA"
  response_time_sla = "10ms"
}
