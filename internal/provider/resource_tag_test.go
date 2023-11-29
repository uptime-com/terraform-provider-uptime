package provider

import (
	"regexp"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestTarResourceImpl(t *testing.T) {
	var (
		_ APIModel                                           = (*TagResourceModel)(nil)
		_ APIModeler[TagResourceModel, upapi.Tag, upapi.Tag] = (*TagResourceModelAdapter)(nil)
		_ API[upapi.Tag, upapi.Tag]                          = (*TagResourceAPI)(nil)
	)
}

func TestAccTagResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_tag/color_hex"),
				ConfigVariables: config.Variables{
					"tag":       config.StringVariable(names[0]),
					"color_hex": config.StringVariable("#ff0000"),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_tag.test", "tag", names[0]),
					resource.TestCheckResourceAttr("uptime_tag.test", "color_hex", "#ff0000"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_tag/color_hex"),
				ConfigVariables: config.Variables{
					"tag":       config.StringVariable(names[0]),
					"color_hex": config.StringVariable("#AA0000"),
				},
				ExpectError: regexp.MustCompile(`Provided configuration value is not a valid hex color`),
			},
		},
	})
}
