package uptime

import (
	"fmt"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCheckSSLCertResource(t *testing.T) {
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
					resource "uptime_check_ssl_cert" "create-update" {
						name               = "%s"
						address            = "example.com"
						contact_groups     = ["nobody"]
						protocol 	       = "https"
						days_before_expiry = 10
					}
				`, tags["create"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_check_ssl_cert.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_check_ssl_cert.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_check_ssl_cert.create-update", "name", tags["create"]),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "uptime_check_ssl_cert" "create-update" {
						name               = "%s"
						address            = "example.com"
						contact_groups     = ["nobody"]
						protocol 	       = "https"
						days_before_expiry = 20
					}
				`, tags["update"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_check_ssl_cert.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_check_ssl_cert.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_check_ssl_cert.create-update", "name", tags["update"]),
				),
			},
		},
	})
}
