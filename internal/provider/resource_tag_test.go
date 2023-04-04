package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testRenderSnippet(t, "resource_tag.tf", 0, nil),
			},
			{
				Config: testRenderSnippet(t, "resource_tag.tf", 1, nil),
			},
		},
	})
}

//func TestAccTagResource(t *testing.T) {
//	api, err := uptimeapi.NewClientWithResponses("https://uptime.com", uptimeapi.WithToken(os.Getenv("UPTIME_TOKEN")))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// There's no guarantee that account used for acceptance testing stays in pristine state at all times. Generate
//	// fairly unique tag names to reduce collision probability.
//	tags := map[string]string{
//		"create":     petname.Generate(2, "-"),
//		"update":     petname.Generate(2, "-"),
//		"import-ok":  petname.Generate(2, "-"),
//		"import-err": petname.Generate(2, "-"),
//	}
//
//	t.Cleanup(func() {
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//
//		obj, err := api.GetServicetaglistWithResponse(ctx, &uptimeapi.GetServicetaglistParams{PageSize: ptr(10000)})
//		if err != nil {
//			t.Fatal(err)
//		}
//		for _, tag := range *obj.JSON200.Results {
//			_, err = api.DeleteServiceTagDetail(ctx, strconv.Itoa(*tag.Pk))
//			if err != nil {
//				t.Fatal(err)
//			}
//		}
//	})
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { _ = testAccAPIClient(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(`
//					resource "uptime_tag" "create-update" {
//						tag       = "%s"
//						color_hex = "#ff0000"
//					}
//				`, tags["create"]),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "id"),
//					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "url"),
//					resource.TestCheckResourceAttr("uptime_tag.create-update", "tag", tags["create"]),
//					resource.TestCheckResourceAttr("uptime_tag.create-update", "color_hex", "#ff0000"),
//				),
//			},
//			{
//				Config: fmt.Sprintf(`
//					resource "uptime_tag" "create-update" {
//						tag       = "%s"
//						color_hex = "#0000ff"
//					}
//				`, tags["update"]),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "id"),
//					resource.TestCheckResourceAttrSet("uptime_tag.create-update", "url"),
//					resource.TestCheckResourceAttr("uptime_tag.create-update", "tag", tags["update"]),
//					resource.TestCheckResourceAttr("uptime_tag.create-update", "color_hex", "#0000ff"),
//				),
//			},
//			{
//				PreConfig: func() {
//					ctx, cancel := context.WithCancel(context.Background())
//					defer cancel()
//
//					param := uptimeapi.GetServicetaglistParams{PageSize: ptr(1)}
//					res, err := api.GetServicetaglistWithResponse(ctx, &param)
//					if err != nil {
//						t.Fatal(err)
//					}
//					obj1 := (*res.JSON200.Results)[0]
//
//					id := strconv.Itoa(*obj1.Pk)
//					obj2 := uptimeapi.CheckTag{
//						Tag:      tags["import-ok"],
//						ColorHex: "#00ff00",
//					}
//
//					_, err = api.PutServiceTagDetailWithResponse(context.Background(), id, obj2)
//					if err != nil {
//						t.Fatal(err)
//					}
//				},
//				ImportState:   true,
//				ImportStateId: tags["import-ok"],
//				ResourceName:  "uptime_tag.import-ok",
//				Config:        `resource "uptime_tag" "import-ok" {}`,
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet("uptime_tag.import-ok", "id"),
//					resource.TestCheckResourceAttrSet("uptime_tag.import-ok", "url"),
//					resource.TestCheckResourceAttr("uptime_tag.import-ok", "tag", tags["import-ok"]),
//					resource.TestCheckResourceAttr("uptime_tag.import-ok", "color_hex", "#00ff00"),
//				),
//			},
//			{
//				ImportState:   true,
//				ImportStateId: tags["import-err"],
//				ResourceName:  "uptime_tag.import-err",
//				Config:        `resource "uptime_tag" "import-err" {}`,
//				ExpectError:   regexp.MustCompile("tag not found"),
//			},
//		},
//	})
//}
