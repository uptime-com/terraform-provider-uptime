package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func LocationsSchemaAttribute(getLocations LocationsSetGetter) schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Default: setdefault.StaticValue(
			types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("US-NY-New York"),
					types.StringValue("US-CA-Los Angeles"),
				},
			),
		),
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Set{
			&locationsPlanModifier{getLocations: getLocations},
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
