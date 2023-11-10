package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
)

func NumRetriesSchemaAttribute(defaultVal int64) schema.Attribute {
	return schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Default:     int64default.StaticInt64(defaultVal),
		Description: "How many times the check should be retried before a location is considered down",
	}
}
