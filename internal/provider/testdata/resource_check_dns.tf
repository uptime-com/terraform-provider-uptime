resource uptime_check_dns test {
  name            = "{{ petname 3 "-" }}"
  contact_groups  = ["nobody"]
  interval        = 5
  locations       = ["US East", "US West"]
  address         = "example.com"
  dns_record_type = "AAAA"
}
// ---
resource uptime_check_dns test {
  name            = "{{ petname 3 "-" }}"
  contact_groups  = ["nobody", "noone"]
  interval        = 10
  locations       = ["Serbia", "Austria"]
  address         = "example.net"
  dns_record_type = "ANY"
}
