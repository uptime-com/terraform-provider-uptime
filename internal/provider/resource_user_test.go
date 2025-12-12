package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/config"
	testResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUserResourceImpl(t *testing.T) {
	var _ resource.Resource = (*UserResource)(nil)
}

func TestAccUserResource(t *testing.T) {
	email := petname.Generate(2, ".") + "@example.com"
	firstName := petname.Generate(1, "-")
	lastName := petname.Generate(1, "-")
	password := petname.Generate(4, "-") + "123"

	testResource.Test(t, testResource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []testResource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_user/basic"),
				ConfigVariables: config.Variables{
					"email":      config.StringVariable(email),
					"first_name": config.StringVariable(firstName),
					"last_name":  config.StringVariable(lastName),
					"password":   config.StringVariable(password),
				},
				Check: testResource.ComposeTestCheckFunc(
					testResource.TestCheckResourceAttr("uptime_user.test", "email", email),
					testResource.TestCheckResourceAttr("uptime_user.test", "first_name", firstName),
					testResource.TestCheckResourceAttr("uptime_user.test", "last_name", lastName),
					testResource.TestCheckResourceAttr("uptime_user.test", "is_api_enabled", "false"),
					testResource.TestCheckResourceAttr("uptime_user.test", "notify_paid_invoices", "false"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_user/basic"),
				ConfigVariables: config.Variables{
					"email":      config.StringVariable(email),
					"first_name": config.StringVariable(firstName),
					"last_name":  config.StringVariable(lastName),
					"password":   config.StringVariable(password),
				},
				ResourceName:            "uptime_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_user/update"),
				ConfigVariables: config.Variables{
					"email":                config.StringVariable(email),
					"first_name":           config.StringVariable(firstName + "-updated"),
					"last_name":            config.StringVariable(lastName + "-updated"),
					"password":             config.StringVariable(password),
					"is_api_enabled":       config.BoolVariable(true),
					"notify_paid_invoices": config.BoolVariable(true),
				},
				Check: testResource.ComposeTestCheckFunc(
					testResource.TestCheckResourceAttr("uptime_user.test", "email", email),
					testResource.TestCheckResourceAttr("uptime_user.test", "first_name", firstName+"-updated"),
					testResource.TestCheckResourceAttr("uptime_user.test", "last_name", lastName+"-updated"),
					testResource.TestCheckResourceAttr("uptime_user.test", "is_api_enabled", "true"),
					testResource.TestCheckResourceAttr("uptime_user.test", "notify_paid_invoices", "true"),
				),
			},
		},
	})
}
