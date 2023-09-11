resource uptime_check_whois test {
  name           = "{{ petname 3 "-" }}"
  address        = "example.com"
  expect_string  = "example.com"
}
// ---
resource uptime_check_whois test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody", "noone"]
  address        = "example.net"
  expect_string  = "example.net"
  threshold      = 10
}
