resource uptime_check_http test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody"]
  locations      = ["US East", "US West"]
  interval       = 5
  address        = "https://example.com"
}
// ---
resource uptime_check_http test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody", "noone"]
  locations      = ["Serbia", "Austria"]
  interval       = 10
  address        = "https://example.net"
}
// ---
resource uptime_check_http test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody", "noone"]
  locations      = ["Serbia", "Austria"]
  interval       = 10
  address        = "https://example.net"
  headers        = {
    "X-My-Header" = ["my-value", "my-other-value"]
  }
}
