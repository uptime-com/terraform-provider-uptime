package provider

import (
	"sort"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDashboardResource(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/_basic"),
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(name),
					"check_name": config.StringVariable(name),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "alerts.include.ignored", "false"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "alerts.include.resolved", "false"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "selected.services.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "selected.services.0", name),
				),
			},
		},
	})
}

func TestAccDashboardResource_RootAttrs(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/_basic"),
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(name),
					"check_name": config.StringVariable(name),
					"ordering":   config.IntegerVariable(10),
					"is_pinned":  config.BoolVariable(false),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_selected.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_selected.0", name),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "is_pinned", "false"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "ordering", "10"),
				),
			},
		},
	})
}

func TestAccDashboardResource_MetricsAttribute(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/_basic"),
				ConfigVariables: config.Variables{
					"name":                   config.StringVariable(name),
					"check_name":             config.StringVariable(name),
					"metrics_show_section":   config.BoolVariable(true),
					"metrics_for_all_checks": config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_selected.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_selected.0", name),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "is_pinned", "false"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "metrics_show_section", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "metrics_for_all_checks", "true"),
				),
			},
		},
	})
}

func TestAccDashboardResource_ServicesAttribute(t *testing.T) {
	t.Parallel()
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/_basic"),
				ConfigVariables: config.Variables{
					"name":                         config.StringVariable(name),
					"check_name":                   config.StringVariable(name),
					"is_pinned":                    config.BoolVariable(false),
					"metrics_show_section":         config.BoolVariable(true),
					"metrics_for_all_checks":       config.BoolVariable(true),
					"services_show_section":        config.BoolVariable(true),
					"services_num_to_show":         config.IntegerVariable(4),
					"services_include_up":          config.BoolVariable(true),
					"services_include_down":        config.BoolVariable(true),
					"services_include_paused":      config.BoolVariable(false),
					"services_include_maintenance": config.BoolVariable(true),
					"services_primary_sort":        config.StringVariable("is_paused,cached_state_is_up"),
					"services_secondary_sort":      config.StringVariable("-cached_last_down_alert_at"),
					"services_show_uptime":         config.BoolVariable(true),
					"services_show_response_time":  config.BoolVariable(true),
					"alerts_show_section":          config.BoolVariable(true),
					"alerts_for_all_checks":        config.BoolVariable(true),
					"alerts_include_ignored":       config.BoolVariable(true),
					"alerts_include_resolved":      config.BoolVariable(true),
					"alerts_num_to_show":           config.IntegerVariable(5),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_selected.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_selected.0", name),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "is_pinned", "false"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "metrics_show_section", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "metrics_for_all_checks", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_show_section", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_num_to_show", "4"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_include_up", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_include_down", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_include_paused", "false"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_include_maintenance", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_primary_sort", "is_paused,cached_state_is_up"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_secondary_sort", "-cached_last_down_alert_at"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_show_uptime", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "services_show_response_time", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "alerts_show_section", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "alerts_for_all_checks", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "alerts_include_ignored", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "alerts_include_resolved", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "alerts_num_to_show", "5"),
				),
			},
		},
	})
}

func TestAccDashboardResource_Tags(t *testing.T) {
	t.Parallel()
	names := []string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	tags := []string{
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
	}
	sort.Strings(tags)
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/tags"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[0]),
				"tags_create": config.SetVariable(
					config.StringVariable(tags[0]),
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
				"tags_use": config.SetVariable(
					config.StringVariable(tags[0]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_dashboard.test", "services_tags.#", "1"),
				resource.TestCheckResourceAttr("uptime_dashboard.test", "services_tags.0", tags[0]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/tags"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(names[1]),
				"tags_create": config.SetVariable(
					config.StringVariable(tags[0]),
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
				"tags_use": config.SetVariable(
					config.StringVariable(tags[1]),
					config.StringVariable(tags[2]),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_dashboard.test", "services_tags.#", "2"),
				resource.TestCheckResourceAttr("uptime_dashboard.test", "services_tags.0", tags[1]),
				resource.TestCheckResourceAttr("uptime_dashboard.test", "services_tags.1", tags[2]),
			),
		},
	}))
}
