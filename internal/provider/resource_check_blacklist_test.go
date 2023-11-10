package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestAccCheckBlacklistResource_impl(t *testing.T) {
	var (
		_ APIModel                                                                   = (*CheckBlacklistResourceModel)(nil)
		_ APIModeler[CheckBlacklistResourceModel, upapi.CheckBlacklist, upapi.Check] = (*CheckBlacklistResourceModelAdapter)(nil)
		_ API[upapi.CheckBlacklist, upapi.Check]                                     = (*CheckBlacklistResourceAPI)(nil)
	)
}

func TestAccCheckBlacklistResource(t *testing.T) {
	t.Parallel()
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"name":    config.StringVariable(names[0]),
					"address": config.StringVariable("example.com"),
				},
				ConfigDirectory: config.StaticDirectory("testdata/resource_check_blacklist/_basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_check_blacklist.test", "name", names[0]),
					resource.TestCheckResourceAttr("uptime_check_blacklist.test", "address", "example.com"),
				),
			},
			{
				ConfigVariables: config.Variables{
					"name":    config.StringVariable(names[1]),
					"address": config.StringVariable("example.net"),
				},
				ConfigDirectory: config.StaticDirectory("testdata/resource_check_blacklist/_basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_check_blacklist.test", "name", names[1]),
					resource.TestCheckResourceAttr("uptime_check_blacklist.test", "address", "example.net"),
				),
			},
		},
	})
}

// TODO: cover attributes
