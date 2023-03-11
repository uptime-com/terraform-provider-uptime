package uptime

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCheckAPIResource(t *testing.T) {
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
resource "uptime_check_api" "create-update" {
	name           = "%s"
	address        = ""
	contact_groups = ["nobody"]
	interval       = 5
	locations      = ["US East", "US West"]
	script         = jsonencode([
		{
			step_def = "C_GET"
			values   = {
				url = "https://example.com/"
			}
		},
	])
}
				`, tags["create"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_check_api.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_check_api.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_check_api.create-update", "name", tags["create"]),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "uptime_check_api" "create-update" {
	name           = "%s"
	address        = ""
	contact_groups = ["nobody", "noone"]
	interval       = 10
	locations      = ["Serbia", "Austria"]
	script         = jsonencode([
		{
			step_def = "C_GET"
			values   = {
				url = "https://example.com/"
			}
		},
	])
}
				`, tags["update"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_check_api.create-update", "name", tags["update"]),
				),
			},
		},
	})
}
