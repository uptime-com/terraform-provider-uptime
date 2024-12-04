package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		},
	})
}
