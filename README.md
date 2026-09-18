# Terraform Provider for Uptime.com

## Installation and usage

Please see the [documentation on the Terraform
Registry](https://registry.terraform.io/providers/uptime-com/uptime/latest/docs).

## Rate Limits

Terraform has a tendency to use many API requests when managing a large group of Uptime.com checks.
Every plan refreshes each check and tag with its own request, so a large account can exhaust its
hourly rate limit, and overlapping runs make this worse. Set `bulk_read = true` on the provider, or
the environment variable `UPTIME_BULK_READ` to any non-empty value, to refresh checks and tags from
their paginated list endpoints instead: a few requests per run regardless of the number of
resources. If this is still not enough, please contact Uptime.com support to request a rate limit
increase.

## Contribution guidelines

See [contribution guidelines](CONTRIBUTING.md).

## Credits

v1 of the provider was created by [Kyle Gentle](https://github.com/kylegentle), with support from Elias
Laham and the Dev Team at Uptime.com.

v2 was created by [Mikhail Lukianchenko](https://github.com/mikluko)

## License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for the full license text.
