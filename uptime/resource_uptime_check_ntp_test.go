package uptime

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCheckNTPResource(t *testing.T) {
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
					resource "uptime_check_ntp" "create-update" {
						name           = "%s"
						address        = "example.com"
						contact_groups = ["nobody"]
						interval       = 5
						locations      = ["US East", "US West"]
					}
				`, tags["create"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_check_ntp.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_check_ntp.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "name", tags["create"]),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "interval", "5"),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "locations.#", "2"),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "locations.0", "US East"),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "locations.1", "US West"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "uptime_check_ntp" "create-update" {
						name           = "%s"
						address        = "example.net"
						contact_groups = ["nobody", "noone"]
						interval       = 10
						locations      = ["Serbia", "US West"]
					}
				`, tags["update"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_check_ntp.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_check_ntp.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "name", tags["update"]),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "locations.#", "2"),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "locations.0", "Serbia"),
					resource.TestCheckResourceAttr("uptime_check_ntp.create-update", "locations.1", "US West"),
				),
			},
		},
	})
}
