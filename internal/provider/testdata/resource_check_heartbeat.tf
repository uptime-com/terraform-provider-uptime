resource uptime_check_heartbeat test {
  name            = "{{ petname 3 "-" }}"
}
// ---
resource uptime_check_heartbeat test {
  name            = "{{ petname 3 "-" }}"
  contact_groups  = ["nobody", "noone"]
  interval        = 10
}
