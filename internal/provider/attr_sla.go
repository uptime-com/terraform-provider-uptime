package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SLASchemaAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "SLA related attributes",
		Attributes: map[string]schema.Attribute{
			"latency": SLALatencySchemaAttribute(),
			"uptime":  SLAUptimeSchemaAttribute(),
		},
		Optional: true,
		Computed: true,
	}
}

type SLAAttribute struct {
	Latency Duration `tfsdk:"latency"`
	Uptime  Decimal  `tfsdk:"uptime"`
}

type SLAAttributeContextAdapter struct{}

func (a SLAAttributeContextAdapter) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"latency": DurationType,
		"uptime":  DecimalType,
	}
}

func (a SLAAttributeContextAdapter) attributeValues(m SLAAttribute) map[string]attr.Value {
	return map[string]attr.Value{
		"latency": m.Latency,
		"uptime":  m.Uptime,
	}
}

func (a SLAAttributeContextAdapter) SLAAttributeContext(ctx context.Context, v types.Object) (*SLAAttribute, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m := *new(SLAAttribute)
	diags := v.As(ctx, &m, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &m, nil
}

func (a SLAAttributeContextAdapter) SLAAttributeValue(m SLAAttribute) types.Object {
	return types.ObjectValueMust(a.attributeTypes(), a.attributeValues(m))
}
