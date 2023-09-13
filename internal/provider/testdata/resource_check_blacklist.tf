resource uptime_check_blacklist test {
  name           = "{{ petname 3 "-" }}"
  address        = "example.com"
}

// ---

resource uptime_check_blacklist test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody", "noone"]
  address        = "example.net"
}
