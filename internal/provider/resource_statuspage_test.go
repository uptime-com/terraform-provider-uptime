package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testRenderSnippet(t, "resource_statuspage.tf", 0, nil),
			},
			{
				Config: testRenderSnippet(t, "resource_statuspage.tf", 1, nil),
			},
			{
				Config: testRenderSnippet(t, "resource_statuspage.tf", 2, nil),
			},
		},
	})
}
