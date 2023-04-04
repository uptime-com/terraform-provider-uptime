resource uptime_check_api test {
  name           = "{{ petname 3 "-" }}"
  contact_groups = ["nobody"]
  interval       = 5
  locations      = ["US East", "US West"]
  script         = jsonencode([
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
}
