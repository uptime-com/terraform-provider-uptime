# Uptime.com Terraform provider changelog

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
