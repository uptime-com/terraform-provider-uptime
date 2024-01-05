package provider

import (
	"regexp"
	"sort"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckTransactionResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "locations.#", "2"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "name", names[1]),
			),
		},
	}))
}

func TestAccCheckTransactionResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckTransactionResource_Interval(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckTransactionResource_Locations(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("US-CA-Los Angeles"),
					config.StringVariable("US-NY-New York"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "locations.0", "US-CA-Los Angeles"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "locations.1", "US-NY-New York"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("Disneyland"),
				),
			},
			ExpectError: regexp.MustCompile(`Invalid value: "Disneyland"`),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("Israel-Tel Aviv"),
					config.StringVariable("Serbia-Belgrade"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "locations.0", "Israel-Tel Aviv"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "locations.1", "Serbia-Belgrade"),
			),
		},
	}))
}

func TestAccCheckTransactionResource_Tags(t *testing.T) {
	name := petname.Generate(3, "-")
	tags := []string{
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
	}
	sort.Strings(tags)
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/tags"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"tags_create": config.SetVariable(
					config.StringVariable(tags[0]),
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
				"tags_use": config.SetVariable(
					config.StringVariable(tags[0]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "tags.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "tags.0", tags[0]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/tags"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"tags_create": config.SetVariable(
					config.StringVariable(tags[0]),
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
				"tags_use": config.SetVariable(
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "tags.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "tags.0", tags[1]),
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "tags.1", tags[2]),
			),
		},
	}))
}

func TestAccCheckTransactionResource_SLA_Uptime(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/sla/uptime"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"sla_uptime": config.StringVariable("0.8"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "sla.uptime", "0.8"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/sla/uptime"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"sla_uptime": config.StringVariable("0.9999"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "sla.uptime", "0.9999"),
			),
		},
	}))
}

func TestAccCheckTransactionResource_SLA_Latency(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("1s"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "sla.latency", "1s"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("60s"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "sla.latency", "60s"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_transaction/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("1000ms"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_transaction.test", "sla.latency", "1000ms"),
			),
		},
	}))
}
