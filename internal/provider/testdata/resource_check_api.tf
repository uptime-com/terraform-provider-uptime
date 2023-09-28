resource uptime_check_api test {
  name   = "{{ petname 3 "-" }}"
  script = jsonencode([
    {
      step_def = "C_GET"
      values   = {
        url = "https://example.com/"
      }
    },
  ])
}
// ---
resource uptime_check_api test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody", "noone"]
  interval       = 10
  locations      = ["Serbia", "Austria"]
  script         = jsonencode([
    {
      step_def = "C_GET"
      values   = {
        url = "https://example.net/"
      }
    },
  ])
  response_time_sla = "500ms"
}
