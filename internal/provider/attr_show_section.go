package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
)

func ShowSectionSchemaAttribute(defaultVal defaults.Bool) schema.Attribute {
	return schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Whether to show the section",
		Default:     defaultVal,
	}
}
