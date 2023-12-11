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
		Optional:    true,
		Computed:    true,
		Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
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
