package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TODO: Extend the test to cover more fields
func TestAccStatusPageResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage/_basic"),
				ConfigVariables: config.Variables{
					"name": config.StringVariable(names[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.basic", "name", names[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage/_basic"),
				ConfigVariables: config.Variables{
					"name": config.StringVariable(names[1]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.basic", "name", names[1]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage/_basic"),
				ConfigVariables: config.Variables{
					"name": config.StringVariable(names[1]),
				},
				ResourceName:            "uptime_statuspage.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_password"},
			},
		},
	})
}

func TestAccStatusPageExtendedResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage/extended"),
				ConfigVariables: config.Variables{
					"name":                         config.StringVariable(names[0]),
					"allow_subscriptions_email":    config.BoolVariable(true),
					"allow_subscriptions_rss":      config.BoolVariable(true),
					"allow_subscriptions_slack":    config.BoolVariable(true),
					"allow_subscriptions_sms":      config.BoolVariable(true),
					"allow_subscriptions_webhook":  config.BoolVariable(true),
					"max_visible_component_days":   config.IntegerVariable(0),
					"hide_empty_tabs_history":      config.BoolVariable(true),
					"theme":                        config.StringVariable("LEGACY"),
					"custom_header_bg_color_hex":   config.StringVariable("#000000"),
					"custom_header_text_color_hex": config.StringVariable("#FFFFFF"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", names[0]),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_email", "true"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_rss", "true"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_slack", "true"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_sms", "true"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_webhook", "true"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "max_visible_component_days", "0"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "hide_empty_tabs_history", "true"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "theme", "LEGACY"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "custom_header_bg_color_hex", "#000000"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "custom_header_text_color_hex", "#FFFFFF"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage/extended"),
				ConfigVariables: config.Variables{
					"name":                         config.StringVariable(names[1]),
					"allow_subscriptions_email":    config.BoolVariable(false),
					"allow_subscriptions_rss":      config.BoolVariable(false),
					"allow_subscriptions_slack":    config.BoolVariable(false),
					"allow_subscriptions_sms":      config.BoolVariable(false),
					"allow_subscriptions_webhook":  config.BoolVariable(false),
					"hide_empty_tabs_history":      config.BoolVariable(false),
					"max_visible_component_days":   config.IntegerVariable(10),
					"theme":                        config.StringVariable("INSPIRE"),
					"custom_header_bg_color_hex":   config.StringVariable("#FFFFFF"),
					"custom_header_text_color_hex": config.StringVariable("#000000"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", names[1]),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_email", "false"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_rss", "false"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_slack", "false"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_sms", "false"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "allow_subscriptions_webhook", "false"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "max_visible_component_days", "10"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "hide_empty_tabs_history", "false"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "theme", "INSPIRE"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "custom_header_bg_color_hex", "#FFFFFF"),
					resource.TestCheckResourceAttr("uptime_statuspage.test", "custom_header_text_color_hex", "#000000"),
				),
			},
		},
	})
}
