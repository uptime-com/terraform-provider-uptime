package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationWavefrontResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	wavefrontURL := [2]string{
		"https://example1.wavefront.com",
		"https://example2.wavefront.com",
	}
	apiToken := [2]string{
		"test-token-1",
		"test-token-2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":          config.StringVariable(names[0]),
				"wavefront_url": config.StringVariable(wavefrontURL[0]),
				"api_token":     config.StringVariable(apiToken[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_wavefront/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_wavefront.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_wavefront.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_wavefront.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_wavefront.test", "wavefront_url", wavefrontURL[0]),
				resource.TestCheckResourceAttr("uptime_integration_wavefront.test", "api_token", apiToken[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":          config.StringVariable(names[1]),
				"wavefront_url": config.StringVariable(wavefrontURL[1]),
				"api_token":     config.StringVariable(apiToken[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_wavefront/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_wavefront.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_wavefront.test", "wavefront_url", wavefrontURL[1]),
				resource.TestCheckResourceAttr("uptime_integration_wavefront.test", "api_token", apiToken[1]),
			),
		},
	}))
}
