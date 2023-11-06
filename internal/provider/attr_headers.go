package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func HeadersSchemaAttribute() schema.Attribute {
	return schema.MapAttribute{
		ElementType: types.ListType{
			ElemType: types.StringType,
		},
		Optional: true,
		Computed: true,
		Default:  mapdefault.StaticValue(types.MapValueMust(types.ListType{ElemType: types.StringType}, map[string]attr.Value{})),
	}
}
