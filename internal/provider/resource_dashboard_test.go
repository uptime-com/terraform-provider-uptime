package provider

import (
	"sort"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDashboardResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "selected.services.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.basic", "selected.services.0", name),
				),
			},
		},
	})
}

func TestAccDashboardResource_Root(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/root"),
				ConfigVariables: config.Variables{
					"name":       config.StringVariable(name),
					"check_name": config.StringVariable(name),
					"is_pinned":  config.BoolVariable(true),
					"ordering":   config.IntegerVariable(11),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.root", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.root", "is_pinned", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.root", "ordering", "11"),
					resource.TestCheckResourceAttr("uptime_dashboard.root", "selected.services.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.root", "selected.services.0", name),
				),
			},
		},
	})
}

func TestAccDashboardResource_Metrics(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/metrics"),
				ConfigVariables: config.Variables{
					"name":                   config.StringVariable(name),
					"check_name":             config.StringVariable(name),
					"metrics_show_section":   config.BoolVariable(true),
					"metrics_for_all_checks": config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.metrics", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.metrics", "metrics.show_section", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.metrics", "metrics.for_all_checks", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.metrics", "selected.services.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.metrics", "selected.services.0", name),
				),
			},
		},
	})
}

func TestAccDashboardResource_Services(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/services"),
				ConfigVariables: config.Variables{
					"name":                         config.StringVariable(name),
					"check_name":                   config.StringVariable(name),
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
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.services", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.show_section", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.num_to_show", "4"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.show.uptime", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.show.response_time", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.include.up", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.include.down", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.include.paused", "false"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.include.maintenance", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.sort.primary", "is_paused,cached_state_is_up"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "services.sort.secondary", "-cached_last_down_alert_at"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "selected.services.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.services", "selected.services.0", name),
				),
			},
		},
	})
}

func TestAccDashboardResource_Alerts(t *testing.T) {
	name := petname.Generate(3, "-")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_dashboard/alerts"),
				ConfigVariables: config.Variables{
					"name":                    config.StringVariable(name),
					"check_name":              config.StringVariable(name),
					"alerts_show_section":     config.BoolVariable(true),
					"alerts_for_all_checks":   config.BoolVariable(true),
					"alerts_num_to_show":      config.IntegerVariable(5),
					"alerts_include_ignored":  config.BoolVariable(true),
					"alerts_include_resolved": config.BoolVariable(true),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "name", name),
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "alerts.show_section", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "alerts.for_all_checks", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "alerts.num_to_show", "5"),
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "alerts.include.ignored", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "alerts.include.resolved", "true"),
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "selected.services.#", "1"),
					resource.TestCheckResourceAttr("uptime_dashboard.alerts", "selected.services.0", name),
				),
			},
		},
	})
}

func TestAccDashboardResource_Tags(t *testing.T) {
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
				resource.TestCheckResourceAttr("uptime_dashboard.tags", "selected.tags.#", "1"),
				resource.TestCheckResourceAttr("uptime_dashboard.tags", "selected.tags.0", tags[0]),
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
				resource.TestCheckResourceAttr("uptime_dashboard.tags", "selected.tags.#", "2"),
				resource.TestCheckResourceAttr("uptime_dashboard.tags", "selected.tags.0", tags[1]),
				resource.TestCheckResourceAttr("uptime_dashboard.tags", "selected.tags.1", tags[2]),
			),
		},
	}))
}
