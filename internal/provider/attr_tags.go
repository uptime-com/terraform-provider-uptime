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
		Description: `List of tags to organize and filter monitoring checks. 
Each account can have up to 3,000 unique tags, with a 100-character limit per tag. 
Tags help categorize resources for filtering in Dashboards, Public Status Pages, and SLA Reports. 
Common use cases include tagging by team ('dev-team', 'ops'), environment ('production', 'staging'), 
or purpose ('api', 'customer-facing'). Defaults to an empty list if not specified.`,
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
