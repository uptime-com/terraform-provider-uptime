package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCheckHeartbeatResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testRenderSnippet(t, "resource_check_heartbeat.tf", 0, nil),
			},
			{
				Config: testRenderSnippet(t, "resource_check_heartbeat.tf", 1, nil),
			},
		},
	})
}
