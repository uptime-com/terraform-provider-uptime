package provider

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCheckAPI(t *testing.T) {
	// There's no guarantee that account used for acceptance testing stays in pristine state at all times. Generate
	// fairly unique tag names to reduce collision probability.
	names := map[string]string{
		"create":     petname.Generate(2, "-"),
		"update":     petname.Generate(2, "-"),
		"import-ok":  petname.Generate(2, "-"),
		"import-err": petname.Generate(2, "-"),
	}

	r.Test(t, r.TestCase{
		PreCheck: func() {
			api := testAccAPIClient(t)
			testAccSetupContactGroup(t, api)
		},
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "uptime_check_api" "test" {
					  name           = "%s"
					  contact_groups = ["void"]
					
					  interval  = 5
					  locations = [
					    "US East",
					    "US West",
					  ]
					
					  script = <<EOF
					    [
					      {
					        "step_def": "C_GET",
					        "values": {
					          "url": "https://httpbin.org/status/200",
					          "headers": {}
					        }
					      },
					      {
					        "step_def": "V_HTTP_STATUS_CODE_IS",
					        "values": {
					          "status_code": "200"
					        }
					      }
					    ]
					  EOF
					}
				`, names["create"]),
				Check: r.ComposeAggregateTestCheckFunc(
					r.TestCheckResourceAttrSet("uptime_check_api.test", "id"),
				),
			},
		},
	})
}
