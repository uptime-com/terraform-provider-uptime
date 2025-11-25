package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationCachetResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	cachetURL := [2]string{
		"https://status.example.com",
		"https://status2.example.com",
	}
	token := [2]string{
		"test-token-123456789", // 20 chars, meets min 20 requirement
		"test-token-234567890", // 20 chars
	}
	component := [2]string{
		"1", // Component ID
		"2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(names[0]),
				"cachet_url": config.StringVariable(cachetURL[0]),
				"token":      config.StringVariable(token[0]),
				"component":  config.StringVariable(component[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_cachet/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_cachet.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_cachet.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_cachet.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_cachet.test", "cachet_url", cachetURL[0]),
				resource.TestCheckResourceAttr("uptime_integration_cachet.test", "token", token[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(names[1]),
				"cachet_url": config.StringVariable(cachetURL[1]),
				"token":      config.StringVariable(token[1]),
				"component":  config.StringVariable(component[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_cachet/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_cachet.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_cachet.test", "cachet_url", cachetURL[1]),
				resource.TestCheckResourceAttr("uptime_integration_cachet.test", "token", token[1]),
			),
		},
	}))
}
