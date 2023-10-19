package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	t.Parallel()
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
					resource.TestCheckResourceAttr("uptime_tag.color_hex", "tag", names[0]),
					resource.TestCheckResourceAttr("uptime_tag.color_hex", "color_hex", "#ff0000"),
				),
			},
		},
	})
}
