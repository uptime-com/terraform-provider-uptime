package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func SLALatencySchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description: `The maximum average response time. Unit is mandatory (e.g. 1500ms or 1.5s or 1s500ms).`,
		Optional:    true,
		Computed:    true,
		Validators: []validator.String{
			slaLatencyValidator{},
		},
		CustomType: DurationType,
	}
}

type slaLatencyValidator struct {
	zoyaDescriber
}

func (d slaLatencyValidator) ValidateString(_ context.Context, rq validator.StringRequest, rs *validator.StringResponse) {
	if rq.ConfigValue.IsNull() {
		return
	}
	if rq.ConfigValue.IsUnknown() {
		return
	}
	_, err := time.ParseDuration(rq.ConfigValue.ValueString())
	if err != nil {
		rs.Diagnostics.AddAttributeError(
			rq.Path,
			"Latency must be a valid duration",
			"Provided configuration value is not a valid duration",
		)
	}
	return
}
