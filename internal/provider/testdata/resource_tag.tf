resource uptime_tag test {
  tag       = "{{ petname 3 "-" }}"
  color_hex = "#ff0000"
}
// ---
resource uptime_tag test {
  tag       = "{{ petname 3 "-" }}"
  color_hex = "#00ff00"
}
