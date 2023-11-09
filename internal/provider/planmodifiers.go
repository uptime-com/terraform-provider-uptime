package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type LocationsSetGetter func(context.Context) (map[string]struct{}, error)

type locationsPlanModifier struct {
	getLocations LocationsSetGetter
}

func (l *locationsPlanModifier) Description(context.Context) string {
	return "Update resource locations with valid values"
}

func (l *locationsPlanModifier) MarkdownDescription(ctx context.Context) string {
	return l.Description(ctx)
}

func (l *locationsPlanModifier) PlanModifySet(
	ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse,
) {
	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		if len(req.PlanValue.Elements()) != 0 {
			locations, err := l.getLocations(ctx)
			if err != nil {
				resp.Diagnostics.AddError("Failed plan modify location set", err.Error())
				return
			}
			for _, el := range req.PlanValue.Elements() {
				if sv, ok := el.(basetypes.StringValue); ok && !sv.IsNull() && !sv.IsUnknown() {
					if _, ok := locations[sv.ValueString()]; !ok {
						resp.Diagnostics.AddError(
							"Failed plan modify location set", "invalid location: "+sv.ValueString())
						return
					}
				}
			}
		}
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.PlanValue
}
