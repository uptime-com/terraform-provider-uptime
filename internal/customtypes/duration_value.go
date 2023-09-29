package customtypes

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.StringValuable = (*Duration)(nil)
)

// Duration represents a time duration.
type Duration struct {
	basetypes.StringValue
}

// Type returns a DurationType.
func (v Duration) Type(_ context.Context) attr.Type {
	return DurationType{}
}

// Equal returns true if the given value is equivalent.
func (v Duration) Equal(x attr.Value) bool {
	o, ok := x.(Duration)
	if !ok {
		return false
	}
	if v.IsNull() && o.IsNull() {
		return true
	}
	if v.IsUnknown() && o.IsUnknown() {
		return true
	}
	dur0, diags0 := v.ValueDuration()
	if diags0.HasError() {
		return false
	}
	dur1, diags1 := o.ValueDuration()
	if diags1.HasError() {
		return false
	}
	return dur0 == dur1
}

// ValueDuration calls time.ParseDuration with the Duration StringValue. A null or unknown value will produce an error diagnostic.
func (v Duration) ValueDuration() (time.Duration, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v.IsNull() {
		diags.Append(diag.NewErrorDiagnostic("Duration ValueDuration Error", "duration string value is null"))
		return time.Duration(0), diags
	}

	if v.IsUnknown() {
		diags.Append(diag.NewErrorDiagnostic("Duration ValueDuration Error", "duration string value is unknown"))
		return time.Duration(0), diags
	}

	dur, err := time.ParseDuration(v.ValueString())
	if err != nil {
		diags.Append(diag.NewErrorDiagnostic("Duration ValueDuration Error", err.Error()))
		return time.Duration(0), diags
	}

	return dur, nil
}

// NewDurationUnknown creates an Duration with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewDurationUnknown() Duration {
	return Duration{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewDurationNull creates an Duration with a null value. Determine whether the value is null via IsNull method.
func NewDurationNull() Duration {
	return Duration{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewDurationValue creates an Duration with a known value.
func NewDurationValue(value string) Duration {
	return Duration{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewDurationPointerValue creates an Duration with a null value if nil or a known value.
func NewDurationPointerValue(value *string) Duration {
	return Duration{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
