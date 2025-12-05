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
that will be notified when alerts are triggered. Defaults to ['Default'] if not specified.
Set to an empty list to disable notifications at this level and rely on parent check group notifications instead.`,
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

type ContactGroupsAttributeAdapter struct {
	SetAttributeAdapter[string]
}

func (a ContactGroupsAttributeAdapter) ContactGroups(v types.Set) *[]string {
	slice := a.Slice(v)
	if slice == nil {
		return nil // null/unknown set - omit field in API call
	}
	return &slice // return pointer (works for empty and non-empty)
}

func (a ContactGroupsAttributeAdapter) ContactGroupsValue(v *[]string) types.Set {
	if v == nil {
		return types.SetNull(types.StringType)
	}
	return a.SliceValue(*v)
}

// ContactGroupsSlice returns []string for API structs that don't use pointer type (integrations).
// Use ContactGroups() for check resources that use *[]string.
func (a ContactGroupsAttributeAdapter) ContactGroupsSlice(v types.Set) []string {
	return a.Slice(v)
}

// ContactGroupsSliceValue accepts []string for API structs that don't use pointer type (integrations).
// Use ContactGroupsValue() for check resources that use *[]string.
func (a ContactGroupsAttributeAdapter) ContactGroupsSliceValue(v []string) types.Set {
	return a.SliceValue(v)
}
