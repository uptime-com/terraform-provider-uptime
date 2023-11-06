package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccCheckHTTPResource(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"address": config.StringVariable("https://example.com"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_check_http.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_check_http.test", "url"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_http.test", "address", "https://example.com"),
			),
			//ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPreRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
				PostApplyPostRefresh: []plancheck.PlanCheck{
					&planCheckNoOp{},
				},
			},
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("https://example.net"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_http.test", "address", "https://example.net"),
			),
		},
	}))
}

func TestAccCheckHTTPResource_ContactGroups(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckHTTPResource_Interval(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckHTTPResource_Locations(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("US-CA-Los Angeles"),
					config.StringVariable("US-NY-New York"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "locations.0", "US-CA-Los Angeles"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "locations.1", "US-NY-New York"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/locations"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"locations": config.ListVariable(
					config.StringVariable("Israel-Tel Aviv"),
					config.StringVariable("Serbia-Belgrade"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "locations.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "locations.0", "Israel-Tel Aviv"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "locations.1", "Serbia-Belgrade"),
			),
		},
	}))
}
