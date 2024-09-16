package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckICMPResource(t *testing.T) {
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
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "address", "example.com"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("example.net"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "address", "example.net"),
			),
		},
	}))
}

func TestAccCheckICMPResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckICMPResource_Interval(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckICMPResource_Locations(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("US-CA-Los Angeles"),
					config.StringVariable("US-NY-New York"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "locations.0", "US-CA-Los Angeles"),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "locations.1", "US-NY-New York"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("Israel-Tel Aviv"),
					config.StringVariable("Serbia-Belgrade"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "locations.0", "Israel-Tel Aviv"),
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "locations.1", "Serbia-Belgrade"),
			),
		},
	}))
}

func TestAccCheckICMPResource_NumRetries(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"num_retries": config.IntegerVariable(3),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "num_retries", "3"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"num_retries": config.IntegerVariable(2),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "num_retries", "2"),
			),
		},
	}))
}

func TestAccCheckICMPResource_Sensitivity(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/sensitivity"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sensitivity": config.IntegerVariable(3),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "sensitivity", "3"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_icmp/sensitivity"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sensitivity": config.IntegerVariable(1),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_icmp.test", "sensitivity", "1"),
			),
		},
	}))
}
