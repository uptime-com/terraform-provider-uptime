package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ColorHexSchemaAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:   true,
		Computed:   true,
		Validators: []validator.String{colorHexValidator{}},
	}
}

type colorHexValidator struct {
	zoyaDescriber
}

func (c colorHexValidator) ValidateString(_ context.Context, rq validator.StringRequest, rs *validator.StringResponse) {
	if rq.ConfigValue.IsNull() || rq.ConfigValue.IsUnknown() {
		return
	}
	v0 := rq.ConfigValue.ValueString()
	v1 := strings.Map(func(r rune) rune {
		if strings.ContainsRune("#0123456789abcdef", r) {
			return r
		}
		return -1
	}, v0)
	if len(v0) != 7 || !strings.HasPrefix(v0, "#") || v0 != v1 {
		rs.Diagnostics.AddAttributeError(
			rq.Path,
			"Provided configuration value is not a valid hex color",
			"Color Hex must be a valid hex color, e.g. #abc0ff (note # prefix and lower case)",
		)
	}
}
