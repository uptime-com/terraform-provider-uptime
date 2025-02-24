package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ContactGroupsSchemaAttribute() schema.SetAttribute {
	return schema.SetAttribute{
		ElementType: types.StringType,
		Description: `List of contact group names to receive notifications. 
Each contact group can contain multiple contacts (email addresses, phone numbers, or integrations) 
that will be notified when alerts are triggered. Defaults to ['Default'] if not specified.`,
		Default: setdefault.StaticValue(
			types.SetValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("Default")},
			),
		),
		Optional: true,
		Computed: true,
	}
}

type ContactGroupsAttribute []string

type ContactGroupsAttributeAdapter struct {
	SetAttributeAdapter[string]
}

func (a ContactGroupsAttributeAdapter) ContactGroups(v types.Set) ContactGroupsAttribute {
	return a.Slice(v)
}

func (a ContactGroupsAttributeAdapter) ContactGroupsValue(v ContactGroupsAttribute) types.Set {
	return a.SliceValue(v)
}
