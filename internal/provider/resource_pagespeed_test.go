package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckPageSpeedResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_pagespeed/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "name", names[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(names[1]),
				"address": config.StringVariable("example.net"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_pagespeed/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "name", names[1]),
			),
		},
	}))
}

func TestAccCheckPageSpeedResource_Config(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_check_pagespeed/config"),
				ConfigVariables: config.Variables{
					"name":                          config.StringVariable(name),
					"pagespeed_config_exclude_urls": config.StringVariable("https://example.com/excluded"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "name", name),
					resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "config.exclude_urls", "https://example.com/excluded"),
				),
			},
		},
	})
}

func TestAccCheckPageSpeedResource_Password(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_check_pagespeed/password"),
				ConfigVariables: config.Variables{
					"name":     config.StringVariable(name),
					"password": config.StringVariable("fakePassword"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "name", name),
					resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "password", "fakePassword"),
				),
			},
		},
	})
}

func TestAccCheckPageSpeedResource_Tags(t *testing.T) {
	name := petname.Generate(3, "-")
	tags := []string{
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_check_pagespeed/tags"),
				ConfigVariables: config.Variables{
					"name": config.StringVariable(name),
					"tags_create": config.SetVariable(
						config.StringVariable(tags[0]),
						config.StringVariable(tags[1]),
					),
					"tags_use": config.SetVariable(config.StringVariable(tags[0])),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "name", name),
					resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("uptime_check_pagespeed.test", "tags.0", tags[0]),
				),
			},
		},
	})
}
