# Uptime.com Terraform provider changelog

## v2.21.0

Enhancements:
* Add Terraform import support for all status page resources:
  * `uptime_statuspage` (simple ID import)
  * `uptime_statuspage_component` (composite ID: `statuspage_id:component_id`)
  * `uptime_statuspage_incident` (composite ID: `statuspage_id:incident_id`)
  * `uptime_statuspage_metric` (composite ID: `statuspage_id:metric_id`)
  * `uptime_statuspage_subscriber` (composite ID: `statuspage_id:subscriber_id`)
  * `uptime_statuspage_user` (composite ID: `statuspage_id:user_id`)
  * `uptime_statuspage_subscription_domain_allow` (composite ID: `statuspage_id:domain_id`)
  * `uptime_statuspage_subscription_domain_block` (composite ID: `statuspage_id:domain_id`)

Dependency Updates:
* Bump github.com/hashicorp/terraform-plugin-framework
* Bump github.com/hashicorp/terraform-plugin-testing

## v2.20.0

New Data Sources:
* uptime_check_groups - Retrieve a list of all check groups

Enhancements:
* `check_group` resource now accepts check IDs (in addition to names) in `config.services`
* `check_maintenance` resource now supports `pause_on_scheduled_maintenance` field

Bug Fixes:
* Fix "Provider produced inconsistent result after apply" error when specifying `check_group` services by numeric ID

Breaking Changes:
* `tag` resource `color_hex` attribute is now required (was optional with computed default)

Documentation:
* Clarify `contact_groups` behavior: defaults to `['Default']` if not specified, set to empty list `[]` to disable notifications

## 2.19.0

This release skipped.

## v2.18.1

Enhancements:
* Allow checks to have empty `contact_groups` (`contact_groups = []`) so alerting can be handled by check groups instead of individual checks

Bug Fixes:
* Fix inverted error check in `UPTIME_RATE_LIMIT` environment variable parsing that caused the parsed value to be ignored

Documentation:
* Add documentation for `check_group` `down_condition` and `uptime_percent_calculation` fields with all valid enum values

## v2.18.0

New Resources:
* uptime_integration_cachet - Cachet status page integration
* uptime_integration_datadog - Datadog monitoring integration
* uptime_integration_geckoboard - Geckoboard dashboard integration
* uptime_integration_jira_servicedesk - JIRA Service Desk integration
* uptime_integration_klipfolio - Klipfolio dashboard integration
* uptime_integration_librato - Librato metrics integration
* uptime_integration_microsoft_teams - Microsoft Teams notification integration
* uptime_integration_pagerduty - PagerDuty incident management integration
* uptime_integration_pushbullet - Pushbullet notification integration
* uptime_integration_pushover - Pushover notification integration
* uptime_integration_slack - Slack notification integration
* uptime_integration_status - Status.io integration
* uptime_integration_statuspage - Statuspage.io integration
* uptime_integration_victorops - VictorOps integration
* uptime_integration_wavefront - Wavefront integration
* uptime_integration_webhook - Custom webhook integration
* uptime_integration_zapier - Zapier automation integration

## v2.17.0

New Data Sources:
* uptime_alerts - Query monitoring alerts and incidents
* uptime_contacts - Query notification recipients
* uptime_dashboards - Query custom monitoring dashboards
* uptime_outages - Query historical downtime events
* uptime_push_notification_profiles - Query mobile device registrations
* uptime_scheduled_reports - Query automated SLA report delivery
* uptime_sla_reports - Query service level agreement tracking
* uptime_statuspage_current_status - Query current status page state
* uptime_statuspage_status_history - Query status page historical data
* uptime_users - Query account team members

New Resource:
* uptime_subaccount - Manage subaccounts

Resource Enhancements:
* uptime_check_sslcert: Added 3 new fields (resolve, ignore_authority_warnings, ignore_sct)

Dependency Updates:
* uptime-client-go/v2: v2.2.1 → v2.4.1
* Various golang.org/x dependencies updated

## v2.16.0

* Add `uptime_user` resource
* Upgrade dependencies

## v2.15.0

* Add `uptime_check_escalations` resource

## v2.14.0

* Add `uptime_statuspages` data source
* Add `uptime_statuspage_components` data source
* Add `uptime_statuspage_incidents` data source
* Add `uptime_statuspage_metrics` data source
* Add `uptime_statuspage_subscribers` data source
* Add `uptime_statuspage_users` data source

## v2.13.0

* Add `uptime_credentials` data source
* Add `uptime_service_variable` resource
* Upgrade dependencies

## v2.12.0

* Add `uptime_check_rdap` resource
* Upgrade dependencies

* Fix IPv4 and IPv6 addresses for locations
* Upgrade dependencies

## v2.11.3

* Fix IPv4 and IPv6 addresses for locations
* Upgrade dependencies

## v2.11.2

* Add `updated_at` to status page incident resource
* Add `max_visible_component_days` to status page component
* Upgrade Go version to 1.24
* Upgrade dependencies

## v2.11.1

* Update Uptime Go client version

## v2.11.0

* Add subaccount ID handling in provider configuration
* Fix dashboard selected items update
* Fix status page component service_id and group_id cleanup
* Fix status page incident test implementation
* Update documentation for complex attributes
* Upgrade dependencies:
  * github.com/hashicorp/terraform-plugin-framework
  * github.com/hashicorp/terraform-plugin-framework-validators
  * github.com/hashicorp/terraform-plugin-go
  * github.com/google/go-cmp from 0.6.0 to 0.7.0

## v2.10.0

* Add `uptime_statuspage_component`, `uptime_statuspage_incident`, `uptime_statuspage_metric`
  `uptime_statuspage_subscriber`, `uptime_statuspage_metric`, `uptime_statuspage_subscription_domain_allow`
  `uptime_statuspage_subscription_domain_block` and `uptime_statuspage_user` resources.
* Add `uptime_credential` resource.
* Add `uptime_integration_opsgenie` resource.
* Remove ICMP check DNS hostname validation for `address` field.
* Upgrade dependencies

## v2.9.0

* Add `uptime_check_maintenance` resource
* Add `sensitivity` to `uptime_check_icmp`, fix #97
* Fix `uptime_percent_calculation` value `AVERAGE` for `uptime_check_group` resource, fix #87
* Upgrade dependencies

## v2.8.0

* Add `uptime_scheduled_report` resource

## v2.7.0

* Add `uptime_sla_report` resource
* Fix TCP check SSL configuration option

## v2.6.0

* Add validation for HTTP check custom port field
* Update documentation
* Upgrade dependencies

## v2.5.2

* Fix goreleaser configuration

## v2.5.1

* Fix `uptime_check_http` resource password property handling
* Upgrade dependencies

## v2.5.0

* Added `uptime_check_udp` resource
* Added `uptime_check_rum` resource
* Added `uptime_check_ssh` resource
* Update TF registry documentation about rate limits
* Upgrade dependencies

## v2.4.3

* Fix default status code for http check
* Upgrade dependencies

## v2.4.2

* Fix transaction check resource
* Upgrade dependencies

## v2.4.1

Changes since v2.4.0

* Upgrade dependencies

## v2.4.0

Changes since v2.3.0

* Added `uptime_check_imap` resource
* Added `uptime_check_pop` resource
* Added `uptime_check_smtp` resource
* Added `uptime_check_transaction` resource

## v2.3.0

Changes since v2.2.0

* Added `uptime_check_webhook` resource
* Added `uptime_check_tcp` resource
* Added `uptime_check_group` resource
* Added `uptime_check_pagespeed` resource

## v2.2.0

Changes since v2.1.2:

* Added `uptime_dashboard` resource

## v2.1.2

Changes since v2.1.1:

* Fixed #36. `notes` attribute now has default value of `"Managed by Terraform"` instead of `""` for all check
  resources. That effectively fixes upgrade from the state created by pre-2.0 version of the provider.

## v2.1.1

Changes since v2.1.0:

* Fixed #35 (merely a dependency version increase, no functional changes)

## v2.1.0

Changes since v2.0.0:

* Added `uptime_contact` resource
* Added `uptime_statuspage` resource
* Added `uptime_check_dns` resource
* Added `uptime_check_icmp` resource
* Added optional SLA parameters to the following resources:
  * `uptime_check_api`
  * `uptime_check_dns`
  * `uptime_check_heartbeat`
  * `uptime_check_http`
  * `uptime_check_icmp`
  * `uptime_check_malware`
  * `uptime_check_ntp`
  * `uptime_check_whois`
* Option `name` used to be optional and failed server-side if not provided. Now it is required for the following
  resources:
  * `uptime_check_api`
  * `uptime_check_blacklist`
  * `uptime_check_dns`
  * `uptime_check_heartbeat`
  * `uptime_check_http`
  * `uptime_check_icmp`
  * `uptime_check_malware`
  * `uptime_check_ntp`
  * `uptime_check_sslcert`
  * `uptime_check_whois`
  * `uptime_contact`
  * `uptime_statuspage`
* The above change for `name` option is functionally backwards compatible since it used to be required by the server
  anyway.
* Option `locations` is now optional and defaults to `["US-NY-New York", "US-CA-Los Angeles"]` for the following
  resources:
  * `uptime_check_api`
  * `uptime_check_http`
  * `uptime_check_ntp`
* Option `locations` now gets validated at Terraform side and fails early instead of being rejected by the server;
* Option `contact_groups` is now optional and defaults to `["Default"]` for the following resources:
  * `uptime_check_api`
  * `uptime_check_blacklist`
  * `uptime_check_dns`
  * `uptime_check_heartbeat`
  * `uptime_check_http`
  * `uptime_check_icmp`
  * `uptime_check_malware`
  * `uptime_check_ntp`
  * `uptime_check_sslcert`
  * `uptime_check_whois`
* Option `theshold` is now optional and defaults to `30` for the following resources:
  * `uptime_check_whois`

## v2.0.0

Initial v2 release. Ported to Terraform Plugin Framework.

Changes since v1.3.4:

* Provider configuration:
  * **BREAKING**: Renamed `rate_limit_ms` to `rate_limit` and changed the unit to requests per second
  * Added `endpoint` argument
  * Added `trace` argument
  * Made `token` optional and configurable from environment variable
* Added `uptime_location` data source
* All resources:
  * `interval` is now optional and determined server-side when omitted
  * Added optional `is_paused`
  * Added optional `num_retries` where applicable
  * Added read-only `url` containing API URL of the resource
* `uptime_check_api`:
  * `address` is removed (not relevant for API checks)
* `uptime_check_dns`:
  * `dns_record_type` is now optional and determined server-side when omitted
* `uptime_check_http`:
  * Added optional `encryption`
  * Added optional `expect_string_type`
  * Added optional `proxy`
  * Added optional `status_code`
  * Added optional `version`
* `uptime_check_malware`:
  * Added read-only `locations`
* `uptime_check_ping`:
  * **BREAKING**: Renamed `ip_version` to `use_ip_version`
* `uptime_check_whois`:
  * **BREAKING**: `expect_string` is now required
  * **BREAKING**: required `days_before_expiry` became optional `threshold`
* **BREAKING**: Renamed `uptime_check_domain_blacklist` resource to `uptime_check_blacklist`
* **BREAKING**: Renamed `uptime_check_ssl_cert` resource to `uptime_check_sslcert`
  * Added nested block for SSL configuration, see docs fro details
* `uptime_tag`:
  * `color_hex` is now optional and determined server-side when omitted
  * `url` read-only property added
* **BREAKING**: Dropped `uptime_integration_opsgenie`
