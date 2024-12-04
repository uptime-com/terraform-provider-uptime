package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageSubscriberResource(t *testing.T) {
	name := petname.Generate(3, "-")
	subscriberTargets := [2]string{
		"test@test.com",
		"https://test.com/webhook",
	}
	subscriberTypes := [2]string{
		"EMAIL",
		"WEBHOOK",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_subscriber/_basic"),
				ConfigVariables: config.Variables{
					"name":              config.StringVariable(name),
					"subscriber_target": config.StringVariable(subscriberTargets[0]),
					"subscriber_type":   config.StringVariable(subscriberTypes[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_subscriber.test", "target", subscriberTargets[0]),
					resource.TestCheckResourceAttr("uptime_statuspage_subscriber.test", "type", subscriberTypes[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_subscriber/_basic"),
				ConfigVariables: config.Variables{
					"name":              config.StringVariable(name),
					"subscriber_target": config.StringVariable(subscriberTargets[1]),
					"subscriber_type":   config.StringVariable(subscriberTypes[1]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_subscriber.test", "target", subscriberTargets[1]),
					resource.TestCheckResourceAttr("uptime_statuspage_subscriber.test", "type", subscriberTypes[1]),
				),
			},
		},
	})
}
