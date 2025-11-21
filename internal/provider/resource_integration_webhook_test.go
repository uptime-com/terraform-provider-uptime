package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationWebhookResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	postbackURL := [2]string{
		"https://example.com/webhook1",
		"https://example.com/webhook2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[0]),
				"postback_url": config.StringVariable(postbackURL[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_webhook/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_webhook.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_webhook.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_webhook.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_webhook.test", "postback_url", postbackURL[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[1]),
				"postback_url": config.StringVariable(postbackURL[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_webhook/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_webhook.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_webhook.test", "postback_url", postbackURL[1]),
			),
		},
	}))
}
