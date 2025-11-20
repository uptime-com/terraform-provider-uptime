package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestSubaccountResourceImpl(t *testing.T) {
	var (
		_ APIModel                                                                             = (*SubaccountResourceModel)(nil)
		_ APIModeler[SubaccountResourceModel, upapi.SubaccountCreateRequest, upapi.Subaccount] = (*SubaccountResourceModelAdapter)(nil)
		_ API[upapi.SubaccountCreateRequest, upapi.Subaccount]                                 = (*SubaccountResourceAPI)(nil)
	)
}

// TestAccSubaccountResource is currently skipped because the subaccounts feature requires
// special permission from Uptime.com support to be enabled.
// API Error: Code=PERMISSION_DENIED Message=Contact support to enable subaccounts feature
// TODO: Re-enable this test once the subaccounts feature is enabled in the test account.
func TestAccSubaccountResource(t *testing.T) {
	t.Skip("Skipping test: Subaccounts feature requires permission from Uptime.com support to be enabled")

	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_subaccount/_basic"),
				ConfigVariables: config.Variables{
					"name": config.StringVariable(names[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_subaccount.test", "id"),
					resource.TestCheckResourceAttrSet("uptime_subaccount.test", "url"),
					resource.TestCheckResourceAttr("uptime_subaccount.test", "name", names[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_subaccount/_basic"),
				ConfigVariables: config.Variables{
					"name": config.StringVariable(names[1]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_subaccount.test", "id"),
					resource.TestCheckResourceAttrSet("uptime_subaccount.test", "url"),
					resource.TestCheckResourceAttr("uptime_subaccount.test", "name", names[1]),
				),
			},
		},
	})
}
