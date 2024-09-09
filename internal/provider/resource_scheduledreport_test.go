package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScheduledReportResource_Basic(t *testing.T) {
	name := petname.Generate(3, "-")
	sla_report_name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
			),
		},
	}))
}

func TestAccScheduledReportResource_RecipientEmails(t *testing.T) {
	name := petname.Generate(3, "-")
	sla_report_name := petname.Generate(3, "-")
	emails := make([]string, 3)
	for i := range emails {
		emails[i] = petname.Generate(2, "@") + ".com"
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"recipient_emails": config.SetVariable(
					config.StringVariable(emails[0]),
					config.StringVariable(emails[1]),
					config.StringVariable(emails[2]),
				),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/recipient_emails"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "recipient_emails.#", "3"),
			),
		},
	}))
}

func TestAccScheduledReportResource_FileType(t *testing.T) {
	name := petname.Generate(3, "-")
	sla_report_name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"file_type":       config.StringVariable("PDF"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/file_type"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "file_type", "PDF"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"file_type":       config.StringVariable("XLS"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/file_type"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "file_type", "XLS"),
			),
		},
	}))
}

func TestAccScheduledReportResource_Recurrence(t *testing.T) {
	name := petname.Generate(3, "-")
	sla_report_name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"recurrence":      config.StringVariable("DAILY"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/recurrence"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "recurrence", "DAILY"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"recurrence":      config.StringVariable("QUARTERLY"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/recurrence"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "recurrence", "QUARTERLY"),
			),
		},
	}))
}

func TestAccScheduledReportResource_OnWeekday(t *testing.T) {
	name := petname.Generate(3, "-")
	sla_report_name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"on_weekday":      config.IntegerVariable(1),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/on_weekday"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "on_weekday", "1"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"on_weekday":      config.IntegerVariable(7),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/on_weekday"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "on_weekday", "7"),
			),
		},
	}))
}

func TestAccScheduledReportResource_AtTime(t *testing.T) {
	name := petname.Generate(3, "-")
	sla_report_name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"at_time":         config.IntegerVariable(0),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/at_time"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "at_time", "0"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"sla_report_name": config.StringVariable(sla_report_name),
				"at_time":         config.IntegerVariable(23),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_scheduledreport/at_time"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "name", name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "sla_report", sla_report_name),
				resource.TestCheckResourceAttr("uptime_scheduled_report.test", "at_time", "23"),
			),
		},
	}))
}
