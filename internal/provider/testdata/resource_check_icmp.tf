resource uptime_check_icmp test {
  name    = "{{ petname 3 "-" }}"
  address = "example.com"
}
// ---
resource uptime_check_icmp test {
  name           = "{{ petname 3 "-" }}"
  locations      = ["Serbia", "Austria"]
  contact_groups = ["nobody", "noone"]
  address        = "example.net"
  num_retries    = 3
  interval       = 10
}
