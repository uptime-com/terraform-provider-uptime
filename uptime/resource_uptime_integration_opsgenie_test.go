package uptime

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIntegrationOpsGenieResource(t *testing.T) {
	// There's no guarantee that account used for acceptance testing stays in pristine state at all times. Generate
	// fairly unique tag names to reduce collision probability.
	tags := map[string]string{
		"create": petname.Generate(2, "-"),
		"update": petname.Generate(2, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoryMap,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "uptime_integration_opsgenie" "create-update" {
						name           = "%s"
						contact_groups = ["nobody"]
						api_endpoint   = "https://api.opsgenie.com/v1"
						api_key        = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
					}
				`, tags["create"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_integration_opsgenie.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_integration_opsgenie.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_integration_opsgenie.create-update", "name", tags["create"]),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "uptime_integration_opsgenie" "create-update" {
						name           = "%s"
						contact_groups = ["nobody", "noone"]
						api_endpoint   = "https://api.opsgenie.com/v2"
						api_key        = "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
					}
				`, tags["update"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_integration_opsgenie.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_integration_opsgenie.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_integration_opsgenie.create-update", "name", tags["update"]),
				),
			},
		},
	})
}
