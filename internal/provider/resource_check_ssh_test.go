package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckSSHResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"address": config.StringVariable("example.com"),
				"port":    config.IntegerVariable(80),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.com"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("example.net"),
				"port":    config.IntegerVariable(80),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.net"),
			),
		},
	}))
}

func TestAccCheckSSHResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckSSHResource_Interval(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"port":     config.IntegerVariable(80),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"port":     config.IntegerVariable(80),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckSSHResource_Locations(t *testing.T) {
	name := petname.Generate(3, "-")
	locs := testAccLocations(t)
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"locations": config.ListVariable(
					config.StringVariable(locs[0]),
					config.StringVariable(locs[1]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.#", "2"),
				resource.TestCheckTypeSetElemAttr("uptime_check_ssh.test", "locations.*", locs[0]),
				resource.TestCheckTypeSetElemAttr("uptime_check_ssh.test", "locations.*", locs[1]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"port": config.IntegerVariable(80),
				"locations": config.ListVariable(
					config.StringVariable(locs[2]),
					config.StringVariable(locs[3]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "locations.#", "2"),
				resource.TestCheckTypeSetElemAttr("uptime_check_ssh.test", "locations.*", locs[2]),
				resource.TestCheckTypeSetElemAttr("uptime_check_ssh.test", "locations.*", locs[3]),
			),
		},
	}))
}

func TestAccCheckSSHResource_NumRetries(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"port":        config.IntegerVariable(80),
				"num_retries": config.IntegerVariable(3),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "num_retries", "3"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/num_retries"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"port":        config.IntegerVariable(80),
				"num_retries": config.IntegerVariable(2),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "num_retries", "2"),
			),
		},
	}))
}

func TestAccCheckSSHResource_Sensitivity(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[0]),
				"address":     config.StringVariable("example.com"),
				"port":        config.IntegerVariable(80),
				"sensitivity": config.IntegerVariable(3),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/sensitivity"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.com"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "sensitivity", "3"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[1]),
				"address":     config.StringVariable("example.net"),
				"port":        config.IntegerVariable(80),
				"sensitivity": config.IntegerVariable(4),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_ssh/sensitivity"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "address", "example.net"),
				resource.TestCheckResourceAttr("uptime_check_ssh.test", "sensitivity", "4"),
			),
		},
	}))
}
