package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationTwitterResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	oauthToken := [2]string{
		"test-oauth-token-1",
		"test-oauth-token-2",
	}
	oauthTokenSecret := [2]string{
		"test-secret-1",
		"test-secret-2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":               config.StringVariable(names[0]),
				"oauth_token":        config.StringVariable(oauthToken[0]),
				"oauth_token_secret": config.StringVariable(oauthTokenSecret[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_twitter/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_twitter.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_twitter.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_twitter.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_twitter.test", "oauth_token", oauthToken[0]),
				resource.TestCheckResourceAttr("uptime_integration_twitter.test", "oauth_token_secret", oauthTokenSecret[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":               config.StringVariable(names[1]),
				"oauth_token":        config.StringVariable(oauthToken[1]),
				"oauth_token_secret": config.StringVariable(oauthTokenSecret[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_twitter/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_twitter.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_twitter.test", "oauth_token", oauthToken[1]),
				resource.TestCheckResourceAttr("uptime_integration_twitter.test", "oauth_token_secret", oauthTokenSecret[1]),
			),
		},
	}))
}
