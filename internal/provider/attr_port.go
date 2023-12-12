package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
)

func PortSchemaAttribute(defaultValue int64) schema.Attribute {
	return schema.Int64Attribute{
		Description: "The port to check",
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(defaultValue),
	}
}

func RequiredPortSchemaAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Description: "The port to check",
		Required:    true,
	}
}
