package provider

import (
	"regexp"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccCheckHTTPResource(t *testing.T) {
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
			// ExpectNonEmptyPlan: true,
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

func TestAccCheckHTTPResource_Headers(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/headers"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"headers": config.MapVariable(map[string]config.Variable{
					"Foo": config.ListVariable(config.StringVariable("Bar")),
				}),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Foo.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Foo.0", "Bar"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/headers"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"headers": config.MapVariable(map[string]config.Variable{
					"Foo": config.ListVariable(config.StringVariable("Bar"), config.StringVariable("Baz")),
					"Qux": config.ListVariable(config.StringVariable("Quux")),
				}),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Foo.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Foo.0", "Bar"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Foo.1", "Baz"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Qux.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_http.test", "headers.Qux.0", "Quux"),
			),
		},
	}))
}

func TestAccCheckHTTPResource_Password(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/password"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"address":  config.StringVariable("https://example.com"),
				"password": config.StringVariable("fakePassword"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "password", "fakePassword"),
			),
		},
	}))
}

func TestAccCheckHTTPResource_PortValidation(t *testing.T) {
	names := []string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/port_validation"),
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"address": config.StringVariable("https://example.com:9383"),
				"port":    config.IntegerVariable(9383),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_http.test", "address", "https://example.com:9383"),
			),
		},
		{
			// basic manifest doesn't contain port definition, so it must fail
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_http/_basic"),
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("https://example.com:9383"),
			},
			ExpectError: regexp.MustCompile("Port value should match"),
		},
	}))
}
