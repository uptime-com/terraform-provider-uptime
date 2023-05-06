package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCheckAPIResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testRenderSnippet(t, "resource_check_api.tf", 0, nil),
			},
			{
				Config: testRenderSnippet(t, "resource_check_api.tf", 1, nil),
			},
		},
	})
}
