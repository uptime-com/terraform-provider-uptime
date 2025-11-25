package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationLibratoResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	email := [2]string{
		"test1@example.com",
		"test2@example.com",
	}
	apiToken := [2]string{
		"test-token-123456789012345678901234567890123456789012345678901234", // 64 chars, meets min 64 requirement
		"test-token-234567890123456789012345678901234567890123456789012345", // 64 chars
	}
	metricName := [2]string{
		"test-metric-1",
		"test-metric-2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[0]),
				"email":       config.StringVariable(email[0]),
				"api_token":   config.StringVariable(apiToken[0]),
				"metric_name": config.StringVariable(metricName[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_librato/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_librato.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_librato.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "email", email[0]),
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "api_token", apiToken[0]),
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "metric_name", metricName[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[1]),
				"email":       config.StringVariable(email[1]),
				"api_token":   config.StringVariable(apiToken[1]),
				"metric_name": config.StringVariable(metricName[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_librato/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "email", email[1]),
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "api_token", apiToken[1]),
				resource.TestCheckResourceAttr("uptime_integration_librato.test", "metric_name", metricName[1]),
			),
		},
	}))
}
