package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// skipIfNoCloudStatusServices skips the test when the account has no
// cloudstatus services or groups provisioned (e.g. the EU test account).
// The SDK does not expose listers for these endpoints yet, so we call them
// directly. Genuine errors (auth, transport, decode, server 5xx) fail the
// test rather than silently skip, so a broken CI setup does not get masked
// as "no data".
func skipIfNoCloudStatusServices(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	token := os.Getenv("UPTIME_TOKEN")
	if token == "" {
		t.Skip("UPTIME_TOKEN must be set for acceptance tests")
	}
	base := os.Getenv("UPTIME_ENDPOINT")
	if base == "" {
		base = "https://uptime.com/api/v1/"
	}
	base = strings.TrimRight(base, "/") + "/"
	client := &http.Client{Timeout: 15 * time.Second}
	countItems := func(path string) (int, error) {
		req, err := http.NewRequest(http.MethodGet, base+path, nil)
		if err != nil {
			return 0, fmt.Errorf("build request %s: %w", path, err)
		}
		req.Header.Set("Authorization", "Token "+token)
		resp, err := client.Do(req)
		if err != nil {
			return 0, fmt.Errorf("GET %s: %w", path, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("read %s body: %w", path, err)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return 0, fmt.Errorf("GET %s: %s: %s", path, resp.Status, string(body))
		}
		var payload struct {
			Count int `json:"count"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return 0, fmt.Errorf("decode %s body: %w", path, err)
		}
		return payload.Count, nil
	}
	svcCount, err := countItems("checks/cloudstatus-services/")
	if err != nil {
		t.Fatalf("cloudstatus precheck: %v", err)
	}
	if svcCount > 0 {
		return
	}
	grpCount, err := countItems("checks/cloudstatus-groups/")
	if err != nil {
		t.Fatalf("cloudstatus precheck: %v", err)
	}
	if grpCount > 0 {
		return
	}
	t.Skip("Skipping: account has no cloudstatus services/groups provisioned")
}

func TestAccCheckCloudStatusResource(t *testing.T) {
	skipIfNoCloudStatusServices(t)
	names := [2]string{
		petname.Generate(3, "-"),
		petname.Generate(3, "-"),
	}
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[0]),
				"service_name": config.StringVariable("Amazon Service"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "name", names[0]),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "service_name", "Amazon Service"),
			),
		},
		{
			ConfigVariables: config.Variables{
				"name":         config.StringVariable(names[1]),
				"service_name": config.StringVariable("100ms API"),
			},
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/_basic"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "name", names[1]),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "service_name", "100ms API"),
			),
		},
	}))
}

func TestAccCheckCloudStatusResource_ContactGroups(t *testing.T) {
	skipIfNoCloudStatusServices(t)
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.#", "1"),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.0", "nobody"),
			),
		},
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/contact_groups"),
			ConfigVariables: config.Variables{
				"name": config.StringVariable(name),
				"contact_groups": config.ListVariable(
					config.StringVariable("Default"),
					config.StringVariable("nobody"),
				),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.#", "2"),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.0", "Default"),
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "contact_groups.1", "nobody"),
			),
		},
	}))
}

func TestAccCheckCloudStatusResource_Group(t *testing.T) {
	skipIfNoCloudStatusServices(t)
	name := petname.Generate(3, "-")
	resource.Test(t, testCaseFromSteps(t, []resource.TestStep{
		{
			ConfigDirectory: config.StaticDirectory("testdata/resource_check_cloudstatus/group"),
			ConfigVariables: config.Variables{
				"name":            config.StringVariable(name),
				"group":           config.IntegerVariable(1),
				"monitoring_type": config.StringVariable("ALL"),
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("uptime_check_cloudstatus.test", "monitoring_type", "ALL"),
			),
		},
	}))
}
