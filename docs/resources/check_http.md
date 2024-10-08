---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "uptime_check_http Resource - terraform-provider-uptime"
subcategory: ""
description: |-
  Monitor a URL for specific status code(s)
---

# uptime_check_http (Resource)

Monitor a URL for specific status code(s)



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String)
- `name` (String)

### Optional

- `contact_groups` (Set of String)
- `encryption` (String) Whether to verify SSL/TLS certificates
- `expect_string` (String)
- `expect_string_type` (String) Valid values for this property are: "STRING" - exact match, "REGEX" - match by regular expression, "INVERSE_REGEX" - fail if the regular expression matches
- `headers` (Map of List of String)
- `include_in_global_metrics` (Boolean) Include this check in uptime/response time calculations for the dashboard and status pages
- `interval` (Number) The interval between checks in minutes
- `is_paused` (Boolean)
- `locations` (Set of String)
- `notes` (String)
- `num_retries` (Number) How many times the check should be retried before a location is considered down
- `password` (String, Sensitive)
- `port` (Number) The `Port` value is mandatory if the address URL contains a custom, non-standard port. It should be set to the same value.
- `proxy` (String)
- `send_string` (String) String to post
- `sensitivity` (Number) How many locations should be down before an alert is sent
- `sla` (Attributes) SLA related attributes (see [below for nested schema](#nestedatt--sla))
- `status_code` (String)
- `tags` (Set of String)
- `threshold` (Number) A timeout alert will be issued if the check takes longer than this many seconds to complete
- `username` (String)
- `version` (Number) Check version to use. Keep default value unless you are absolutely sure you need to change it

### Read-Only

- `id` (Number) The ID of this resource.
- `url` (String)

<a id="nestedatt--sla"></a>
### Nested Schema for `sla`

Optional:

- `latency` (String) The maximum average response time. Unit is mandatory (e.g. 1500ms or 1.5s or 1s500ms).
- `uptime` (String) The minimum uptime percentage. \nMust be a fraction with exactly 4 decimal places (e.g. 0.9995 for 99.95% uptime)


