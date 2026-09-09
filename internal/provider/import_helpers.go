package provider

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ParseCompositeID parses an import ID in the format "<parentAttr>:resource_id" and
// returns both IDs as int64 values. parentAttr names the resource's own parent attribute
// so the error text tells the user the exact format that resource expects.
func ParseCompositeID(id string, parentAttr string) (parentID int64, resourceID int64, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected format '%s:resource_id', got '%s'", parentAttr, id)
	}

	parentID, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid %s '%s': %w", parentAttr, parts[0], err)
	}
	if err := checkIDRange(parentID, parentAttr, parts[0]); err != nil {
		return 0, 0, err
	}

	resourceID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid resource_id '%s': %w", parts[1], err)
	}
	if err := checkIDRange(resourceID, "resource_id", parts[1]); err != nil {
		return 0, 0, err
	}

	return parentID, resourceID, nil
}

// checkIDRange rejects import IDs that cannot be used as a primary key.
//
// Database keys start at 1, so a non-positive half is always a typo. It matters most for
// the resource id: 0 is the value that reads back as "gone", so an unvalidated 0 would
// import a resource that every later plan proposes recreating.
//
// The upper bound is the platform int, not int64, because upapi.PrimaryKey is an int and
// the provider ships 386 and arm builds - on those a larger id would silently truncate and
// address a different record.
func checkIDRange(id int64, attr string, raw string) error {
	if id < 1 || id > math.MaxInt {
		return fmt.Errorf("invalid %s '%s': must be a positive integer no larger than %d", attr, raw, math.MaxInt)
	}
	return nil
}

// ImportStateSimpleIDFor returns an import handler for resources with a single numeric
// key, writing it into idAttr. Most resources key on "id", but a resource whose schema
// names its key differently must name that attribute here, or the write finds no such
// attribute.
func ImportStateSimpleIDFor(
	idAttr string,
) func(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		id, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID",
				fmt.Sprintf("expected numeric %s, got '%s': %s", idAttr, req.ID, err.Error()))
			return
		}
		if err := checkIDRange(id, idAttr, req.ID); err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", err.Error())
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(idAttr), id)...)
	}
}

// ImportStateSimpleID handles import for resources with a simple numeric ID.
// It parses the import ID string and sets it as an int64 "id" attribute.
func ImportStateSimpleID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStateSimpleIDFor("id")(ctx, req, resp)
}

// ImportStateCompositeIDFor returns an import handler for child resources whose
// composite ID is "<parent>:<resource_id>", writing the parent half into parentAttr.
// The parent attribute must be imported explicitly whenever the API does not return
// it, otherwise Read cannot repopulate it and a Required parent would plan a spurious
// change right after import.
func ImportStateCompositeIDFor(
	parentAttr string,
) func(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	return func(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
		parentID, resourceID, err := ParseCompositeID(req.ID, parentAttr)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", err.Error())
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(parentAttr), parentID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceID)...)
	}
}

// ImportStateCompositeID handles import for child resources with composite IDs.
// It parses the import ID in format "statuspage_id:resource_id" and sets both attributes.
func ImportStateCompositeID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ImportStateCompositeIDFor("statuspage_id")(ctx, req, resp)
}
