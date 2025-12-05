package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ServicesSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Description: `List of check IDs to be included in the group (specified as strings, e.g., ["5581024"]).
A group can contain up to 200 individual checks of any type (except other group checks).
Checks can be part of multiple groups simultaneously. Defaults to an empty list if not specified.`,
		Optional: true,
		Computed: true,
		Default:  setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
	}
}

type ServicesAttribute []string

type ServicesAttributeAdapter struct {
	SetAttributeAdapter[string]
}

func (a ServicesAttributeAdapter) Tags(v types.Set) ServicesAttribute {
	return a.Slice(v)
}

func (a ServicesAttributeAdapter) TagsValue(v ServicesAttribute) types.Set {
	return a.SliceValue(v)
}
