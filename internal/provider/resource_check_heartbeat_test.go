package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckHeartbeatResource(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "name", names[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "name", names[1]),
			),
		},
	}))
}

func TestAccCheckHeartbeatResource_ContactGroups(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckHeartbeatResource_Interval(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(5),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "interval", "5"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/interval"),
			ConfigVariables: config.Variables{
				"name":     config.StringVariable(name),
				"interval": config.IntegerVariable(10),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "interval", "10"),
			),
		},
	}))
}

func TestAccCheckHeartbeatResource_SLA_Uptime(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/sla/uptime"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"sla_uptime": config.StringVariable("0.8"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "sla.uptime", "0.8"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/sla/uptime"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"sla_uptime": config.StringVariable("0.9999"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "sla.uptime", "0.9999"),
			),
		},
	}))
}

func TestAccCheckHeartbeatResource_SLA_Latency(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("1s"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "sla.latency", "1s"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("60s"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "sla.latency", "60s"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_heartbeat/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("1000ms"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_heartbeat.test", "sla.latency", "1000ms"),
			),
		},
	}))
}
