package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckWebhookResource(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_webhook/_basic"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_webhook.test", "name", name),
				resource.TestCheckResourceAttrSet("uptime_check_webhook.test", "webhook_url"),
			),
		},
	}))
}

func TestAccCheckWebhookResource_ContactGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_webhook/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
					config.StringVariable("noone"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_webhook.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_webhook.test", "contact_groups.0", "nobody"),
				resource.TestCheckResourceAttr("uptime_check_webhook.test", "contact_groups.1", "noone"),
			),
		},
	}))
}

func TestAccCheckWebhookResource_SLA_Uptime(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_webhook/sla/uptime"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"sla_uptime": config.StringVariable("0.9999"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_webhook.test", "sla.uptime", "0.9999"),
			),
		},
	}))
}

func TestAccCheckWebhookResource_SLA_Latency(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_webhook/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("10s"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_webhook.test", "sla.latency", "10s"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_webhook/sla/latency"),
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(name),
				"sla_latency": config.StringVariable("500ms"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_webhook.test", "sla.latency", "500ms"),
			),
		},
	}))
}
