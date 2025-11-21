package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationStatuspageResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	apiKey := [2]string{
		"test-api-key-1",
		"test-api-key-2",
	}
	page := [2]string{
		"test-page-id-1",
		"test-page-id-2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[0]),
				"api_key": config.StringVariable(apiKey[0]),
				"page":    config.StringVariable(page[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_statuspage/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_statuspage.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_statuspage.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_statuspage.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_statuspage.test", "api_key", apiKey[0]),
				resource.TestCheckResourceAttr("uptime_integration_statuspage.test", "page", page[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"api_key": config.StringVariable(apiKey[1]),
				"page":    config.StringVariable(page[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_statuspage/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_statuspage.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_statuspage.test", "api_key", apiKey[1]),
				resource.TestCheckResourceAttr("uptime_integration_statuspage.test", "page", page[1]),
			),
		},
	}))
}
