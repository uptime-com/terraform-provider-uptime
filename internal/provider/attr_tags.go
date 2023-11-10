package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TagsSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
		Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
	}
}

type TagsAttribute []string

type TagsAttributeAdapter struct {
	SetAttributeAdapter[string]
}

func (a TagsAttributeAdapter) Tags(v types.Set) TagsAttribute {
	return a.Slice(v)
}

func (a TagsAttributeAdapter) TagsValue(v TagsAttribute) types.Set {
	return a.SliceValue(v)
}
