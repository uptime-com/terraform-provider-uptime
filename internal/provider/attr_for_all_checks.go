package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ForAllChecksSchemaAttribute() schema.Attribute {
	return schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Whether to show block for all checks",
	}
}
