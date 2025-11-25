package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationVictoropsResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	serviceKey := [2]string{
		"test-service-key-1234567890123456789", // 37 chars (min 36)
		"test-service-key-2345678901234567890", // 37 chars
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[0]),
				"service_key": config.StringVariable(serviceKey[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_victorops/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_victorops.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_victorops.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_victorops.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_victorops.test", "service_key", serviceKey[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(names[1]),
				"service_key": config.StringVariable(serviceKey[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_victorops/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_victorops.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_victorops.test", "service_key", serviceKey[1]),
			),
		},
	}))
}
