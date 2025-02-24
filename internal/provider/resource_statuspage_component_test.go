package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageComponentResource(t *testing.T) {
	name := petname.Generate(3, "-")
	componentNames := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[1]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[1]),
				),
			},
		},
	})
}

func TestAccStatusPageComponentResource_Clean(t *testing.T) {
	name := petname.Generate(3, "-")
	componentNames := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/clean_step1"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/clean_step2"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[1]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[1]),
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "description", ""),
					resource.TestCheckNoResourceAttr("uptime_statuspage_component.test", "group_id"),
					resource.TestCheckNoResourceAttr("uptime_statuspage_component.test", "service_id"),
				),
			},
		},
	})
}
