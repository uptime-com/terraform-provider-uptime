package provider

import (
	"sort"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSLAReportResource_Basic(t *testing.T) {
	name := petname.Generate(3, "-")
	check_name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"check_name": config.StringVariable(check_name),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_slareport/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_sla_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "services_selected.#", "1"),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "services_selected.0.name", check_name),
			),
		},
	}))
}

func TestAccSLAReportResource_ReportingGroups(t *testing.T) {
	name := petname.Generate(3, "-")
	check_name := petname.Generate(3, "-")
	reporting_group_name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":                 config.StringVariable(name),
				"check_name":           config.StringVariable(check_name),
				"reporting_group_name": config.StringVariable(reporting_group_name),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_slareport/reporting_groups"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_sla_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "reporting_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "reporting_groups.0.name", reporting_group_name),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "reporting_groups.0.group_services.#", "1"),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "reporting_groups.0.group_services.0", check_name),
			),
		},
	}))
}

func TestAccSLAReportResource_Tags(t *testing.T) {
	name := petname.Generate(3, "-")
	check_name := petname.Generate(3, "-")
	tags := []string{
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
		petname.Generate(2, "-"),
	}
	sort.Strings(tags)
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_slareport/tags"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"check_name": config.StringVariable(check_name),
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
				resource.TestCheckResourceAttr("uptime_sla_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "services_tags.#", "1"),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "services_tags.0", tags[0]),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_slareport/tags"),
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(name),
				"check_name": config.StringVariable(check_name),
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
				resource.TestCheckResourceAttr("uptime_sla_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "services_tags.#", "2"),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "services_tags.0", tags[1]),
				resource.TestCheckResourceAttr("uptime_sla_report.test", "services_tags.1", tags[2]),
			),
		},
	}))
}
