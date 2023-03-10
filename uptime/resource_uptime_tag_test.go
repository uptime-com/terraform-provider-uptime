package uptime

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"github.com/uptime-com/uptime-client-go"
)

func TestAccTagResource(t *testing.T) {
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
					resource "uptime_tag" "create-update" {
						tag       = "%s"
						color_hex = "#ff0000"
					}
				`, tags["create"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_tag.create-update", "tag", tags["create"]),
					resource.TestCheckResourceAttr("uptime_tag.create-update", "color_hex", "#ff0000"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "uptime_tag" "create-update" {
						tag       = "%s"
						color_hex = "#0000ff"
					}
				`, tags["update"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "id"),
					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "url"),
					resource.TestCheckResourceAttr("uptime_tag.create-update", "tag", tags["update"]),
					resource.TestCheckResourceAttr("uptime_tag.create-update", "color_hex", "#0000ff"),
				),
			},
			{
				ImportState: true,
				ImportStateIdFunc: func(_ *terraform.State) (string, error) {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()

					api := testAccAPIClient(t)
					obj, res, err := api.Tags.Create(ctx, &uptime.Tag{
						Tag:      tags["import-ok"],
						ColorHex: "#00ff00",
					})
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, res.StatusCode)

					t.Cleanup(func() {
						_, _ = api.Tags.Delete(context.Background(), obj.PK)
					})
					return strconv.Itoa(obj.PK), nil
				},
				ResourceName: "uptime_tag.import-ok",
				Config:       `resource "uptime_tag" "import-ok" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("uptime_tag.import-ok", "id"),
					resource.TestCheckResourceAttrSet("uptime_tag.import-ok", "url"),
					resource.TestCheckResourceAttr("uptime_tag.import-ok", "tag", tags["import-ok"]),
					resource.TestCheckResourceAttr("uptime_tag.import-ok", "color_hex", "#00ff00"),
				),
			},
			{
				ImportState:   true,
				ImportStateId: tags["import-err"],
				ResourceName:  "uptime_tag.import-err",
				Config:        `resource "uptime_tag" "import-err" {}`,
				ExpectError:   regexp.MustCompile("NOT_FOUND"),
			},
		},
	})
}
