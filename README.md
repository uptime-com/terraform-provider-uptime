# Terraform Provider for Uptime.com
## Requirements
* Terraform v0.12.0 or higher
* Go v1.12 or higher

## Installation
### Downloading the provider
First, install the provider to your local machine:
```
go get -u -v github.com/uptime-com/terraform-provider-uptime
```

### Installing
In order for Terraform to use terraform-provider-uptime, it needs to be linked to the plugin directory. Example commands for an OS X Darwin machine:

```
mkdir -p ~/.terraform.d/plugins/darwin_amd64
ln -s ~/go/bin/terraform-provider-uptime ~/.terraform.d/plugins/darwin_amd64/terraform-provider-uptime
```

For Linux machines, follow the OS X process, replacing `darwin` with `linux`.

For a Windows machine, in PowerShell:
```
New-Item %APPDATA%\terraform.d\plugins\windows_amd64 -Type 'directory' -Force
cmd /c mklink /d $env:GOPATH\bin\terraform-provider-uptime %APPDATA%\terraform.d\plugins\windows_amd64\terraform-provider-uptime
```

## Resources
(Section is incomplete)

#### uptime\_check\_ntp
Example:

```
resource "uptime_check_ntp" "google" {
    name = "Google Public NTP"
    address = "time.google.com"
    contact_groups = ["Default", "NTP"]
    interval = 1
    locations = ["US-East", "GBR"]
    tags = ["terraform"]
}
```

Required attributes:

* **address**, *string*: address of the server under test

* **contact_groups**, *list(string)*: contact groups to alert

* **interval**, *number*: time interval between checks, in minutes

* **locations**, *list(string)*: probe server locations

Optional attributes:

* **name**, *string*: human-readable/friendly name

* **tags**, *list(string)*: tags to attach to the check

* **notes**, *string*: arbitrary notes for check

* **include_in_global_metrics**, *bool*: whether to include this check in global uptime metrics

* **ip_version**, *string, limited to "IPV4" or "IPV6"*: IP version to use

* **sensitivity**, *number*: number of probe servers that should detect a failure before an alert is triggered

* **threshold**, *number*: timeout threshold for server response, in seconds

* **port**, *number*: port where service is running

## Credits
terraform-provider-uptime was originally created by [Kyle Gentle](https://github.com/kylegentle), with support from Elias Laham and the Dev Team at Uptime.com.
