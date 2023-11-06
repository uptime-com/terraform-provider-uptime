package provider

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ScriptSchemaAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		CustomType:  RawJsonType{},
		Description: `The script to run. Must be valid JSON.`,
		Optional:    true,
		Computed:    true,
		Validators: []validator.String{
			scriptValidator{},
		},
	}
}

type scriptValidator struct {
	zoyaDescriber
}

func (s scriptValidator) ValidateString(_ context.Context, rq validator.StringRequest, rs *validator.StringResponse) {
	if rq.ConfigValue.IsNull() || rq.ConfigValue.IsUnknown() {
		return
	}
	if !json.Valid(bytes.NewBufferString(rq.ConfigValue.String()).Bytes()) {
		rs.Diagnostics.AddAttributeError(
			rq.Path,
			"Script must be valid JSON",
			"Provided configuration value is not valid JSON",
		)
	}
}
