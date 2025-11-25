package provider

import (
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationJiraServicedeskResource(t *testing.T) {
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	apiEmail := [2]string{
		"test1@example.com",
		"test2@example.com",
	}
	apiToken := [2]string{
		"test-token-1",
		"test-token-2",
	}
	jiraSubdomain := [2]string{
		"https://example1.atlassian.net", // Valid URL format
		"https://example2.atlassian.net", // Valid URL format
	}
	projectKey := [2]string{
		"TEST1",
		"TEST2",
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":           config.StringVariable(names[0]),
				"api_email":      config.StringVariable(apiEmail[0]),
				"api_token":      config.StringVariable(apiToken[0]),
				"jira_subdomain": config.StringVariable(jiraSubdomain[0]),
				"project_key":    config.StringVariable(projectKey[0]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_jira_servicedesk/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("uptime_integration_jira_servicedesk.test", "id"),
				resource.TestCheckResourceAttrSet("uptime_integration_jira_servicedesk.test", "url"),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "api_email", apiEmail[0]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "api_token", apiToken[0]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "jira_subdomain", jiraSubdomain[0]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "project_key", projectKey[0]),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":           config.StringVariable(names[1]),
				"api_email":      config.StringVariable(apiEmail[1]),
				"api_token":      config.StringVariable(apiToken[1]),
				"jira_subdomain": config.StringVariable(jiraSubdomain[1]),
				"project_key":    config.StringVariable(projectKey[1]),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_integration_jira_servicedesk/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "api_email", apiEmail[1]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "api_token", apiToken[1]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "jira_subdomain", jiraSubdomain[1]),
				resource.TestCheckResourceAttr("uptime_integration_jira_servicedesk.test", "project_key", projectKey[1]),
			),
		},
	}))
}
