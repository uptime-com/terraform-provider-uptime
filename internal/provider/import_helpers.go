package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ParseCompositeID parses an import ID in the format "parent_id:resource_id"
// and returns both IDs as int64 values.
func ParseCompositeID(id string) (parentID int64, resourceID int64, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected format 'statuspage_id:resource_id', got '%s'", id)
	}

	parentID, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid statuspage_id '%s': %w", parts[0], err)
	}

	resourceID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid resource_id '%s': %w", parts[1], err)
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

// ImportStateCompositeID handles import for child resources with composite IDs.
// It parses the import ID in format "statuspage_id:resource_id" and sets both attributes.
func ImportStateCompositeID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	statusPageID, resourceID, err := ParseCompositeID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("statuspage_id"), statusPageID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceID)...)
}
