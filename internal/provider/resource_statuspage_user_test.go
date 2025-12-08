package provider

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccStatusPageUserResource(t *testing.T) {
	t.Skip("Skipping test as API is broken: GET/DELETE methods")
	name := petname.Generate(3, "-")
	emails := [2]string{
		petname.Generate(3, "-") + "@test.com",
		petname.Generate(3, "-") + "@test.com",
	}
	firstNames := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	lastNames := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_user/_basic"),
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(name),
					"email":      config.StringVariable(emails[0]),
					"first_name": config.StringVariable(firstNames[0]),
					"last_name":  config.StringVariable(lastNames[0]),
					"is_active":  config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "email", emails[0]),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "first_name", firstNames[0]),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "last_name", lastNames[0]),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "is_active", "true"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_user/_basic"),
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(name),
					"email":      config.StringVariable(emails[1]),
					"first_name": config.StringVariable(firstNames[1]),
					"last_name":  config.StringVariable(lastNames[1]),
					"is_active":  config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "email", emails[1]),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "first_name", firstNames[1]),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "last_name", lastNames[1]),
					resource.TestCheckResourceAttr("uptime_statuspage_user.test", "is_active", "false"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_user/_basic"),
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(name),
					"email":      config.StringVariable(emails[1]),
					"first_name": config.StringVariable(firstNames[1]),
					"last_name":  config.StringVariable(lastNames[1]),
					"is_active":  config.BoolVariable(false),
				},
				ResourceName:      "uptime_statuspage_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					statusPageRS := s.RootModule().Resources["uptime_statuspage.test"]
					userRS := s.RootModule().Resources["uptime_statuspage_user.test"]
					if statusPageRS == nil || userRS == nil {
						return "", fmt.Errorf("resources not found in state")
					}
					return fmt.Sprintf("%s:%s", statusPageRS.Primary.Attributes["id"], userRS.Primary.Attributes["id"]), nil
				},
			},
		},
	})
}
