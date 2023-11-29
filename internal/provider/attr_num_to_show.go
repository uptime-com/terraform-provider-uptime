package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
)

func NumToShowSchemaAttribute(defaultValue defaults.Int64) schema.Attribute {
	return schema.Int64Attribute{
		Description: "The number of entities to show",
		Optional:    true,
		Computed:    true,
		Default:     defaultValue,
	}
}
