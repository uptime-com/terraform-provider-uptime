package uptime

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestHeaderConversion(t *testing.T) {
	hm := make(map[string]interface{})
	hm["Foo"] = "bar"
	hm["Baz"] = "bat"

	require.Contains(t, headersMapToString(hm), "Foo: bar")
	require.Contains(t, headersMapToString(hm), "Baz: bat")
}

func TestAccCheckHttpResource(t *testing.T) {
	// There's no guarantee that account used for acceptance testing stays in pristine state at all times. Generate
	// fairly unique tag names to reduce collision probability.
	tags := map[string]string{
		"create":     petname.Generate(2, "-"),
		"update":     petname.Generate(2, "-"),
		"import-ok":  petname.Generate(2, "-"),
		"import-err": petname.Generate(2, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoryMap,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "uptime_check_http" "create-update" {
						name           = "%s"
						address        = "https://example.com"
						contact_groups = ["nobody"]
						interval       = 5
						locations      = ["US East", "US West"]
					}
				`, tags["create"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_check_http.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_check_http.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "name", tags["create"]),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "contact_groups.#", "1"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "contact_groups.0", "nobody"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "interval", "5"),
					resource.TestCheckResourceAttrSet("uptime_check_http.create-update", "num_retries"),
					func(state *terraform.State) error {
						return nil
					},
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "uptime_check_http" "create-update" {
						name           = "%s"
						address        = "https://example.net"
						contact_groups = ["nobody", "noone"]
						interval       = 10
						locations      = ["Serbia"]
						num_retries    = 3
					}
				`, tags["update"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_check_http.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_check_http.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "name", tags["update"]),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "contact_groups.#", "2"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "contact_groups.0", "nobody"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "contact_groups.1", "noone"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "interval", "10"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "locations.#", "1"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "locations.0", "Serbia"),
					resource.TestCheckResourceAttr("uptime_check_http.create-update", "num_retries", "3"),
				),
			},
		},
	})
}
