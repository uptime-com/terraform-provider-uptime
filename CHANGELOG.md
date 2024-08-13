# Uptime.com Terraform provider changelog

## v2.6.0
- Add validation for HTTP check custom port field
- Update documentation
- Upgrade dependencies

## v2.5.2

- Fix goreleaser configuration

## v2.5.1

- Fix `uptime_check_http` resource password property handling
- Upgrade dependencies

## v2.5.0

- Added `uptime_check_udp` resource
- Added `uptime_check_rum` resource
- Added `uptime_check_ssh` resource
- Update TF registry documentation about rate limits
- Upgrade dependencies

## v2.4.3

- Fix default status code for http check
- Upgrade dependencies

## v2.4.2

- Fix transaction check resource
- Upgrade dependencies

## v2.4.1

Changes since v2.4.0

- Upgrade dependencies

## v2.4.0

Changes since v2.3.0

- Added `uptime_check_imap` resource
- Added `uptime_check_pop` resource
- Added `uptime_check_smtp` resource
- Added `uptime_check_transaction` resource

## v2.3.0

Changes since v2.2.0

- Added `uptime_check_webhook` resource
- Added `uptime_check_tcp` resource
- Added `uptime_check_group` resource
- Added `uptime_check_pagespeed` resource

## v2.2.0

Changes since v2.1.2:

- Added `uptime_dashboard` resource

## v2.1.2

Changes since v2.1.1:

- Fixed #36. `notes` attribute now has default value of `"Managed by Terraform"` instead of `""` for all check
  resources. That effectively fixes upgrade from the state created by pre-2.0 version of the provider.

## v2.1.1

Changes since v2.1.0:

- Fixed #35 (merely a dependency version increase, no functional changes)

## v2.1.0

Changes since v2.0.0:

- Added `uptime_contact` resource
- Added `uptime_statuspage` resource
- Added `uptime_check_dns` resource
- Added `uptime_check_icmp` resource
- Added optional SLA parameters to the following resources:
  - `uptime_check_api`
  - `uptime_check_dns`
  - `uptime_check_heartbeat`
  - `uptime_check_http`
  - `uptime_check_icmp`
  - `uptime_check_malware`
  - `uptime_check_ntp`
  - `uptime_check_whois`
- Option `name` used to be optional and failed server-side if not provided. Now it is required for the following
  resources:
  - `uptime_check_api`
  - `uptime_check_blacklist`
  - `uptime_check_dns`
  - `uptime_check_heartbeat`
  - `uptime_check_http`
  - `uptime_check_icmp`
  - `uptime_check_malware`
  - `uptime_check_ntp`
  - `uptime_check_sslcert`
  - `uptime_check_whois`
  - `uptime_contact`
  - `uptime_statuspage`
- The above change for `name` option is functionally backwards compatible since it used to be required by the server
  anyway.
- Option `locations` is now optional and defaults to `["US-NY-New York", "US-CA-Los Angeles"]` for the following
  resources:
  - `uptime_check_api`
  - `uptime_check_http`
  - `uptime_check_ntp`
- Option `locations` now gets validated at Terraform side and fails early instead of being rejected by the server;
- Option `contact_groups` is now optional and defaults to `["Default"]` for the following resources:
  - `uptime_check_api`
  - `uptime_check_blacklist`
  - `uptime_check_dns`
  - `uptime_check_heartbeat`
  - `uptime_check_http`
  - `uptime_check_icmp`
  - `uptime_check_malware`
  - `uptime_check_ntp`
  - `uptime_check_sslcert`
  - `uptime_check_whois`
- Option `theshold` is now optional and defaults to `30` for the following resources:
  - `uptime_check_whois`

## v2.0.0

Initial v2 release. Ported to Terraform Plugin Framework.

Changes since v1.3.4:

- Provider configuration:
  - **BREAKING**: Renamed `rate_limit_ms` to `rate_limit` and changed the unit to requests per second
  - Added `endpoint` argument
  - Added `trace` argument
  - Made `token` optional and configurable from environment variable
- Added `uptime_location` data source
- All resources:
  - `interval` is now optional and determined server-side when omitted
  - Added optional `is_paused`
  - Added optional `num_retries` where applicable
  - Added read-only `url` containing API URL of the resource
- `uptime_check_api`:
  - `address` is removed (not relevant for API checks)
- `uptime_check_dns`:
  - `dns_record_type` is now optional and determined server-side when omitted
- `uptime_check_http`:
  - Added optional `encryption`
  - Added optional `expect_string_type`
  - Added optional `proxy`
  - Added optional `status_code`
  - Added optional `version`
- `uptime_check_malware`:
  - Added read-only `locations`
- `uptime_check_ping`:
  - **BREAKING**: Renamed `ip_version` to `use_ip_version`
- `uptime_check_whois`:
  - **BREAKING**: `expect_string` is now required
  - **BREAKING**: required `days_before_expiry` became optional `threshold`
- **BREAKING**: Renamed `uptime_check_domain_blacklist` resource to `uptime_check_blacklist`
- **BREAKING**: Renamed `uptime_check_ssl_cert` resource to `uptime_check_sslcert`
  - Added nested block for SSL configuration, see docs fro details
- `uptime_tag`:
  - `color_hex` is now optional and determined server-side when omitted
  - `url` read-only property added
- **BREAKING**: Dropped `uptime_integration_opsgenie`
