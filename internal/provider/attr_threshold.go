package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ThresholdSchemaAttribute(defaultVal int64) schema.Int64Attribute {
	return ThresholdDescriptionSchemaAttribute(defaultVal, "A timeout alert will be issued if the check takes longer than this many seconds to complete")
}

func ThresholdDescriptionSchemaAttribute(defaultVal int64, desc string) schema.Int64Attribute {
	return schema.Int64Attribute{
		CustomType:  types.Int64Type,
		Default:     int64default.StaticInt64(defaultVal),
		Optional:    true,
		Computed:    true,
		Description: desc,
	}
}
