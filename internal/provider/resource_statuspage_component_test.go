package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestAccStatusPageComponentResource(t *testing.T) {
	name := petname.Generate(3, "-")
	componentNames := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[1]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[1]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[1]),
				},
				ResourceName:      "uptime_statuspage_component.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					statusPageRS := s.RootModule().Resources["uptime_statuspage.test"]
					componentRS := s.RootModule().Resources["uptime_statuspage_component.test"]
					if statusPageRS == nil || componentRS == nil {
						return "", fmt.Errorf("resources not found in state")
					}
					return fmt.Sprintf("%s:%s", statusPageRS.Primary.Attributes["id"], componentRS.Primary.Attributes["id"]), nil
				},
			},
		},
	})
}

// TestAccStatusPageComponentResource_RefreshAfterDelete verifies that when a
// component is deleted out-of-band (e.g. its status page or group was
// cascade-deleted), a refresh drops it from state instead of failing the run
// (SYS-1180).
func TestAccStatusPageComponentResource_RefreshAfterDelete(t *testing.T) {
	name := petname.Generate(3, "-")
	componentName := petname.Generate(3, "-")

	var statusPageID, componentID int64
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentName),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("uptime_statuspage.test", "id", func(v string) error {
						id, err := strconv.ParseInt(v, 10, 64)
						statusPageID = id
						return err
					}),
					resource.TestCheckResourceAttrWith("uptime_statuspage_component.test", "id", func(v string) error {
						id, err := strconv.ParseInt(v, 10, 64)
						componentID = id
						return err
					}),
				),
			},
			{
				PreConfig: func() {
					api := testAccAPIClient(t)
					err := api.StatusPages().
						Components(upapi.PrimaryKey(statusPageID)).
						Delete(context.Background(), upapi.PrimaryKey(componentID))
					require.NoError(t, err, "out-of-band component delete failed")
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccStatusPageComponentResource_SortingWeight(t *testing.T) {
	name := petname.Generate(3, "-")
	componentName := petname.Generate(3, "-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/sorting_weight"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentName),
					"sorting_weight": config.IntegerVariable(10),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentName),
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "sorting_weight", "10"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/sorting_weight"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentName),
					"sorting_weight": config.IntegerVariable(500),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "sorting_weight", "500"),
				),
			},
		},
	})
}

func TestAccStatusPageComponentResource_Clean(t *testing.T) {
	name := petname.Generate(3, "-")
	componentNames := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/clean_step1"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[0]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[0]),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_component/clean_step2"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"component_name": config.StringVariable(componentNames[1]),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "name", componentNames[1]),
					resource.TestCheckResourceAttr("uptime_statuspage_component.test", "description", ""),
					resource.TestCheckNoResourceAttr("uptime_statuspage_component.test", "group_id"),
					resource.TestCheckNoResourceAttr("uptime_statuspage_component.test", "service_id"),
				),
			},
		},
	})
}
