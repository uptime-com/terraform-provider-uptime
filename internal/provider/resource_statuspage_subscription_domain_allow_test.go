package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageSubscriptionDomainAllowResource(t *testing.T) {
	name := petname.Generate(3, "-")
	domains := [2]string{
		"test1.com",
		"test2.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_subscription_domain_allow/_basic"),
				ConfigVariables: config.Variables{
					"name":   config.StringVariable(name),
					"domain": config.StringVariable(domains[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_subscription_domain_allow.test", "domain", domains[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_subscription_domain_allow/_basic"),
				ConfigVariables: config.Variables{
					"name":   config.StringVariable(name),
					"domain": config.StringVariable(domains[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_subscription_domain_allow.test", "domain", domains[0]),
				),
			},
		},
	})
}
