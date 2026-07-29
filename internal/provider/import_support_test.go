package provider

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

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

	out := make(map[string]struct{ implemented, documented bool })
	for _, newResource := range p.Resources(ctx) {
		r := newResource()
		metaResp := fwresource.MetadataResponse{}
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "uptime"}, &metaResp)

		_, implemented := r.(fwresource.ResourceWithImportState)
		_, err := os.Stat(filepath.Join("..", "..", "examples", "resources", metaResp.TypeName, "import.sh"))
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

// TestSimpleImportWritesKeyAttribute runs the import handler of the resources that key on a
// single numeric attribute against their real schema. uptime_check_maintenance is the reason
// this is worth asserting: it keys on check_id and has no id attribute, so the generic "id"
// handler would write to an attribute that does not exist and fail only at import time.
func TestSimpleImportWritesKeyAttribute(t *testing.T) {
	ctx := context.Background()
	p := &providerImpl{}

	for _, tc := range []struct {
		newResource func() fwresource.Resource
		typeName    string
		keyAttr     string
	}{
		{func() fwresource.Resource { return NewCredentialResource(ctx, p) }, "uptime_credential", "id"},
		{func() fwresource.Resource { return NewCheckMaintenanceResource(ctx, p) }, "uptime_check_maintenance", "check_id"},
	} {
		t.Run(tc.typeName, func(t *testing.T) {
			r := tc.newResource()
			schemaResp := fwresource.SchemaResponse{}
			r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
			resp := &fwresource.ImportStateResponse{State: tfsdk.State{
				Schema: schemaResp.Schema,
				Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
			}}

			importer, ok := r.(fwresource.ResourceWithImportState)
			if !ok {
				t.Fatalf("%s does not implement resource.ResourceWithImportState", tc.typeName)
			}
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
