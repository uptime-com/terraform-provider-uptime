package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func SLAUptimeSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description: `The minimum uptime percentage. \n` +
			`Must be a fraction with exactly 4 decimal places (e.g. 0.9995 for 99.95% uptime)`,
		Optional:   true,
		Computed:   true,
		CustomType: DecimalType,
	}
}
