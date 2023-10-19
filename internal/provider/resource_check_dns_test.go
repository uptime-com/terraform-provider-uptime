package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckDNSResource(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"address": config.StringVariable("example.com"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "address", "example.com"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("example.net"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "address", "example.net"),
			),
		},
	}))
}

func TestAccCheckDNSResource_ContactGroups(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckDNSResource_DNSRecordType(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/dns_record_type"),
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"dns_record_type": config.StringVariable("A"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "dns_record_type", "A"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/dns_record_type"),
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"dns_record_type": config.StringVariable("AAAA"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "dns_record_type", "AAAA"),
			),
		},
	}))
}

func TestAccCheckDNSResource_Interval(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckDNSResource_Locations(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("US East"),
					config.StringVariable("US West"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "locations.0", "US East"),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "locations.1", "US West"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("Austria"),
					config.StringVariable("Serbia"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "locations.0", "Austria"),
				resource.TestCheckResourceAttr("uptime_check_dns.test", "locations.1", "Serbia"),
			),
		},
	}))
}

func TestAccCheckDNSResource_ResponseTimeSLA(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/response_time_sla"),
			ConfigVariables: config.Variables{
				"name":              config.StringVariable(name),
				"response_time_sla": config.StringVariable("100ms"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "response_time_sla", "100ms"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_dns/response_time_sla"),
			ConfigVariables: config.Variables{
				"name":              config.StringVariable(name),
				"response_time_sla": config.StringVariable("200ms"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_dns.test", "response_time_sla", "200ms"),
			),
		},
	}))
}
