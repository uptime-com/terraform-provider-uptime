package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageIncidentBasicResource(t *testing.T) {
	name := petname.Generate(3, "-")
	incidentName := petname.Generate(3, "-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_incident/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"incident_name":  config.StringVariable(incidentName),
					"incident_state": config.StringVariable("investigating"),
					"starts_at":      config.StringVariable("2025-01-28T00:00:00Z"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "name", incidentName),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "incident_type", "INCIDENT"),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "starts_at", "2025-01-28T00:00:00Z"),
				),
			},
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_incident/_basic"),
				ConfigVariables: config.Variables{
					"name":           config.StringVariable(name),
					"incident_name":  config.StringVariable(incidentName),
					"incident_state": config.StringVariable("monitoring"),
					"starts_at":      config.StringVariable("2025-01-28T10:10:00Z"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "updates.#", "1"),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "updates.0.incident_state", "monitoring"),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "starts_at", "2025-01-28T10:10:00Z"),
				),
			},
		},
	})
}

func TestAccStatusPageIncidentAffectedComponentsResource(t *testing.T) {
	name := petname.Generate(3, "-")
	incidentName := petname.Generate(3, "-")
	checkName := petname.Generate(3, "-")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { _ = testAccAPIClient(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/resource_statuspage_incident/affected_components"),
				ConfigVariables: config.Variables{
					"name":                      config.StringVariable(name),
					"incident_name":             config.StringVariable(incidentName),
					"check_name":                config.StringVariable(checkName),
					"incident_state":            config.StringVariable("investigating"),
					"incident_component_status": config.StringVariable("major-outage"),
					"starts_at":                 config.StringVariable("2025-01-28T00:00:00Z"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptime_statuspage.test", "name", name),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "name", incidentName),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "incident_type", "INCIDENT"),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "affected_components.#", "1"),
					resource.TestCheckResourceAttr(
						"uptime_statuspage_incident.test", "affected_components.0.status", "major-outage",
					),
					resource.TestCheckResourceAttr("uptime_statuspage_incident.test", "starts_at", "2025-01-28T00:00:00Z"),
				),
			},
		},
	})
}
