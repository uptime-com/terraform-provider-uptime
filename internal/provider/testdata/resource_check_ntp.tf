resource uptime_check_ntp test {
  name    = "{{ petname 3 "-" }}"
  address = "example.com"
}
// ---
resource uptime_check_ntp test {
  name              = "{{ petname 3 "-" }}"
  contact_groups    = ["nobody", "noone"]
  locations         = ["Serbia", "Austria"]
  interval          = 10
  address           = "example.net"
  response_time_sla = "100ms"
}
