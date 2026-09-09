package provider

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// exampleResourcesDir holds the per-resource example trees, relative to this package.
// go test runs the binary in the package directory, so the relative base is stable.
const exampleResourcesDir = "../../examples/resources"

// resourceImportSupport reports, for every registered resource, whether it implements
// resource.ResourceWithImportState and whether the docs advertise import for it.
//
// tfplugindocs renders the "## Import" section of docs/resources/<name>.md purely from the
// presence of examples/resources/uptime_<name>/import.sh, so that file is what a user reads
// as "this resource is importable" - nothing checks it against the code.
func resourceImportSupport(t *testing.T) map[string]struct{ implemented, documented bool } {
	t.Helper()
	ctx := context.Background()
	p := &providerImpl{}

	// Without this the checks degrade silently: if the examples tree moves, every resource
	// reads as undocumented, "documented && !implemented" is never true, and the guard that
	// matters passes green while checking nothing.
	if _, err := os.Stat(exampleResourcesDir); err != nil {
		t.Fatalf("cannot locate %s: %v - the import checks derive documented support from "+
			"that tree and would silently pass if it moved", exampleResourcesDir, err)
	}

	out := make(map[string]struct{ implemented, documented bool })
	for _, newResource := range p.Resources(ctx) {
		r := newResource()
		metaResp := fwresource.MetadataResponse{}
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "uptime"}, &metaResp)

		_, implemented := r.(fwresource.ResourceWithImportState)
		_, err := os.Stat(filepath.Join(exampleResourcesDir, metaResp.TypeName, "import.sh"))
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("stat import example for %s: %v", metaResp.TypeName, err)
		}
		out[metaResp.TypeName] = struct{ implemented, documented bool }{implemented, err == nil}
	}
	if len(out) == 0 {
		t.Fatal("no resources registered - the provider wiring changed and this check is now blind")
	}
	return out
}

// TestResourceImportDocsMatchCode guards the SYS-1303 failure mode: uptime_service_variable
// carried an ImportState method on a type the framework never consults, so it looked
// importable in review while `terraform import` answered "Resource Import Not Implemented".
// The customer only found out at their terminal, following instructions we had written.
//
// Documentation is the promise users act on, so it must not claim import the code does not
// implement. The reverse is a documentation gap, not a broken promise, and is reported
// separately.
func TestResourceImportDocsMatchCode(t *testing.T) {
	for name, support := range resourceImportSupport(t) {
		if support.documented && !support.implemented {
			t.Errorf("%s: docs advertise terraform import but the resource does not implement "+
				"resource.ResourceWithImportState - construct it with NewImportableAPIResource", name)
		}
	}
}

// importsNeedingLiveAPI lists resources whose import handler calls the API, so it cannot be
// exercised against a zero-value provider. uptime_service_variable cross-checks that the
// service_id half of its composite ID really owns the variable (SYS-1303), which is two live
// reads. Its handler is covered by TestServiceVariableImportState with a stubbed client.
var importsNeedingLiveAPI = map[string]bool{"uptime_service_variable": true}

// documentedImportID returns the ID argument from a resource's import.sh, i.e. the exact
// string the docs tell a user to run.
func documentedImportID(t *testing.T, typeName string) (string, bool) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(exampleResourcesDir, typeName, "import.sh"))
	if err != nil {
		return "", false
	}
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[0] == "terraform" && fields[1] == "import" {
			return fields[len(fields)-1], true
		}
	}
	t.Errorf("%s: import.sh contains no `terraform import` command", typeName)
	return "", false
}

// TestDocumentedImportIDIsAccepted feeds each resource the ID its own documentation tells
// users to run, through its real handler and real schema.
//
// The docs-vs-code check above only proves import exists. This proves the documented call
// actually works - it is what catches a resource being switched between a simple and a
// composite key while its example keeps the old shape, which reproduces the SYS-1303
// experience exactly: the docs say run X, the terminal rejects X.
func TestDocumentedImportIDIsAccepted(t *testing.T) {
	ctx := context.Background()
	p := &providerImpl{}

	for _, newResource := range p.Resources(ctx) {
		r := newResource()
		metaResp := fwresource.MetadataResponse{}
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "uptime"}, &metaResp)

		importer, ok := r.(fwresource.ResourceWithImportState)
		if !ok || importsNeedingLiveAPI[metaResp.TypeName] {
			continue
		}
		id, ok := documentedImportID(t, metaResp.TypeName)
		if !ok {
			continue
		}

		t.Run(metaResp.TypeName, func(t *testing.T) {
			schemaResp := fwresource.SchemaResponse{}
			r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
			resp := &fwresource.ImportStateResponse{State: tfsdk.State{
				Schema: schemaResp.Schema,
				Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
			}}

			importer.ImportState(ctx, fwresource.ImportStateRequest{ID: id}, resp)
			if resp.Diagnostics.HasError() {
				t.Errorf("documented import ID %q rejected: %v", id, resp.Diagnostics.Errors())
			}
		})
	}
}

// TestSimpleImportWritesKeyAttribute pins that the key lands in the attribute the resource
// actually uses. The negative cases pin the diagnostics, which must name the resource's own
// key attribute - the whole reason ImportStateSimpleIDFor takes one.
func TestSimpleImportWritesKeyAttribute(t *testing.T) {
	ctx := context.Background()
	p := &providerImpl{}

	for _, tc := range []struct {
		newResource func() fwresource.Resource
		typeName    string
		keyAttr     string
	}{
		{func() fwresource.Resource { return NewCredentialResource(ctx, p) }, "uptime_credential", "id"},
	} {
		newState := func() *fwresource.ImportStateResponse {
			schemaResp := fwresource.SchemaResponse{}
			tc.newResource().Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
			return &fwresource.ImportStateResponse{State: tfsdk.State{
				Schema: schemaResp.Schema,
				Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
			}}
		}
		importer := tc.newResource().(fwresource.ResourceWithImportState)

		t.Run(tc.typeName+"/writes the key", func(t *testing.T) {
			resp := newState()
			importer.ImportState(ctx, fwresource.ImportStateRequest{ID: "123"}, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("import failed: %v", resp.Diagnostics.Errors())
			}
			var got types.Int64
			if diags := resp.State.GetAttribute(ctx, path.Root(tc.keyAttr), &got); diags.HasError() {
				t.Fatalf("reading %s: %v", tc.keyAttr, diags.Errors())
			}
			if got.ValueInt64() != 123 {
				t.Errorf("%s = %d, want 123", tc.keyAttr, got.ValueInt64())
			}
		})

		// "0" matters most: it is the id that reads back as "gone", so importing it would
		// leave a resource every later plan proposes recreating.
		for _, bad := range []string{"0", "-1", "abc", "", "99999999999999999999"} {
			t.Run(tc.typeName+"/rejects "+strconv.Quote(bad), func(t *testing.T) {
				resp := newState()
				importer.ImportState(ctx, fwresource.ImportStateRequest{ID: bad}, resp)
				if !resp.Diagnostics.HasError() {
					t.Fatalf("import ID %q was accepted", bad)
				}
				detail := resp.Diagnostics.Errors()[0].Detail()
				if !strings.Contains(detail, "numeric "+tc.keyAttr+",") &&
					!strings.Contains(detail, "invalid "+tc.keyAttr+" '") {
					t.Errorf("diagnostic should name %s, got: %s", tc.keyAttr, detail)
				}
			})
		}
	}
}

// TestResourceImportImplementedIsDocumented reports importable resources with no import
// example: `terraform import` works but the docs never say so, and tfplugindocs renders no
// Import section. Undocumented capability is the milder direction of the same drift.
func TestResourceImportImplementedIsDocumented(t *testing.T) {
	for name, support := range resourceImportSupport(t) {
		if support.implemented && !support.documented {
			t.Errorf("%s: implements terraform import but has no examples/resources/%s/import.sh, "+
				"so the generated docs omit the Import section", name, name)
		}
	}
}
