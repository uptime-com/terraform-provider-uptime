package provider

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccStatusPageMetricResource(t *testing.T) {
	name := petname.Generate(3, "-")
	metricName := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_metric/_basic"),
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(name),
					"metric_name": config.StringVariable(metricName[0]),
					"is_visible":  config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_metric.test", "name", metricName[0]),
					resource.TestCheckResourceAttr("uptime_statuspage_metric.test", "is_visible", "false"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_metric/_basic"),
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(name),
					"metric_name": config.StringVariable(metricName[1]),
					"is_visible":  config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_metric.test", "name", metricName[1]),
					resource.TestCheckResourceAttr("uptime_statuspage_metric.test", "is_visible", "true"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_metric/_basic"),
				ConfigVariables: config.Variables{
					"name":        config.StringVariable(name),
					"metric_name": config.StringVariable(metricName[1]),
					"is_visible":  config.BoolVariable(true),
				},
				ResourceName:      "uptime_statuspage_metric.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					statusPageRS := s.RootModule().Resources["uptime_statuspage.test"]
					metricRS := s.RootModule().Resources["uptime_statuspage_metric.test"]
					if statusPageRS == nil || metricRS == nil {
						return "", fmt.Errorf("resources not found in state")
					}
					return fmt.Sprintf("%s:%s", statusPageRS.Primary.Attributes["id"], metricRS.Primary.Attributes["id"]), nil
				},
			},
		},
	})
}
