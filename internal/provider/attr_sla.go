package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SLASchemaAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "SLA related attributes. When omitted, the server-managed SLA is left untouched; " +
			"to clear it, set explicit zero values (`uptime = \"0\"`, `latency = \"0s\"`).",
		Attributes: map[string]schema.Attribute{
			"latency": SLALatencySchemaAttribute(),
			"uptime":  SLAUptimeSchemaAttribute(),
		},
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Object{
			SLADefaultPreserveState(),
		},
	}
}

// SLADefaultPreserveState keeps the prior state value for `sla` when the user
// did not set it in config. Without this, omitting the sla block plans a
// zero-valued object, which the API then writes back as msp_*_sla=0, wiping a
// server-managed SLA and churning `~ sla` on every apply.
func SLADefaultPreserveState() planmodifier.Object {
	return &slaDefaultPreserveState{}
}

type slaDefaultPreserveState struct{}

func (s *slaDefaultPreserveState) Description(context.Context) string {
	return "Preserve the prior state value when the config does not set sla."
}

func (s *slaDefaultPreserveState) MarkdownDescription(ctx context.Context) string {
	return s.Description(ctx)
}

func (s *slaDefaultPreserveState) PlanModifyObject(_ context.Context, rq planmodifier.ObjectRequest, rs *planmodifier.ObjectResponse) {
	if !rq.ConfigValue.IsNull() {
		return
	}
	if rq.StateValue.IsNull() || rq.StateValue.IsUnknown() {
		return
	}
	rs.PlanValue = rq.StateValue
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
