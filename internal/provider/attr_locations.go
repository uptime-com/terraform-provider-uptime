package provider

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func LocationsSchemaAttribute(l LocationsGetter) schema.SetAttribute {
	return LocationsSchemaAttributeWithDefaults(l, "United Kingdom-London", "Netherlands-Amsterdam")
}

func LocationsSchemaAttributeWithDefaults(l LocationsGetter, defaults ...string) schema.SetAttribute {
	defaultValues := make([]attr.Value, len(defaults))
	for i := range defaults {
		defaultValues[i] = types.StringValue(defaults[i])
	}
	return schema.SetAttribute{
		ElementType: types.StringType,
		Default: setdefault.StaticValue(
			types.SetValueMust(types.StringType, defaultValues),
		),
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Set{
			LocationsPlanModifier(l),
		},
	}
}

func LocationsReadOnlySchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}
}

type LocationsAttribute []string

type LocationsAttributeAdapter struct {
	SetAttributeAdapter[string]
}

func (a LocationsAttributeAdapter) Locations(v types.Set) LocationsAttribute {
	return a.SetAttributeAdapter.Slice(v)
}

func (a LocationsAttributeAdapter) LocationsValue(v LocationsAttribute) types.Set {
	return a.SetAttributeAdapter.SliceValue(v)
}

type LocationsGetter interface {
	GetLocations(context.Context) (map[string]struct{}, error)
}

func LocationsPlanModifier(l LocationsGetter) planmodifier.Set {
	return &locationsPlanModifier{LocationsGetter: l}
}

type locationsPlanModifier struct {
	LocationsGetter
}

func (l *locationsPlanModifier) Description(context.Context) string {
	return "Update resource locations with valid values"
}

func (l *locationsPlanModifier) MarkdownDescription(ctx context.Context) string {
	return l.Description(ctx)
}

func (l *locationsPlanModifier) PlanModifySet(ctx context.Context, rq planmodifier.SetRequest, rs *planmodifier.SetResponse) {
	// Do nothing if there is a known planned value.
	if rq.PlanValue.IsUnknown() {
		return
	}
	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if rq.ConfigValue.IsUnknown() {
		return
	}
	// Do nothing if planned value is empty.
	if len(rq.PlanValue.Elements()) == 0 {
		return
	}

	locMap, err := l.GetLocations(ctx)
	if err != nil {
		rs.Diagnostics.AddError("Failed to get valid locations set", err.Error())
		return
	}

	for _, el := range rq.PlanValue.Elements() {
		sv, ok := el.(types.String)
		if !ok {
			rs.Diagnostics.AddError("Location set element is not a string", fmt.Sprintf("Actual value is %T: %v ", el, el))
			return
		}
		if sv.IsNull() || sv.IsUnknown() {
			rs.Diagnostics.AddError("Location set element is null or unknown", "")
			return
		}
		if _, ok = locMap[sv.ValueString()]; !ok {
			locSlc := make([]string, 0, len(locMap))
			for k := range locMap {
				locSlc = append(locSlc, k)
			}
			sort.Strings(locSlc)

			b := *new(strings.Builder)
			b.WriteString("Invalid value: " + sv.String() + "\n\nValid values: \n")
			for i := range locSlc {
				b.WriteString("  - " + locSlc[i] + "\n")
			}
			rs.Diagnostics.AddError("Location is not valid", b.String())
			return
		}
	}

	rs.PlanValue = rq.PlanValue
}
