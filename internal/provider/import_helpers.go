package provider

import (
	"context"
	"fmt"
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
	if parentID < 1 {
		return 0, 0, fmt.Errorf("invalid %s '%s': must be a positive integer", parentAttr, parts[0])
	}

	resourceID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid resource_id '%s': %w", parts[1], err)
	}
	// Database primary keys start at 1, so a non-positive half is always a typo. Rejecting
	// it here matters most for the resource id: 0 is the value that reads back as "gone",
	// so an unvalidated 0 would import a resource that every later plan proposes recreating.
	if resourceID < 1 {
		return 0, 0, fmt.Errorf("invalid resource_id '%s': must be a positive integer", parts[1])
	}

	return parentID, resourceID, nil
}

// ImportStateSimpleID handles import for resources with a simple numeric ID.
// It parses the import ID string and sets it as an int64 "id" attribute.
func ImportStateSimpleID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID",
			fmt.Sprintf("expected numeric ID, got '%s': %s", req.ID, err.Error()))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
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
