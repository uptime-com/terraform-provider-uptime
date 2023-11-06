package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/shopspring/decimal"
)

// DurationType is an attribute type that represents a time duration.
var DurationType = durationType{}

type durationType struct{}

func (t durationType) TerraformType(context.Context) tftypes.Type {
	return tftypes.String
}

// StringValue returns a human readable string of the type name.
func (t durationType) String() string {
	return "provider.Duration"
}

// ValueType returns the Value type.
func (t durationType) ValueType(context.Context) attr.Value {
	return durationValue{}
}

// Equal returns true if the given type is equivalent.
func (t durationType) Equal(o attr.Type) bool {
	_, ok := o.(durationType)
	return ok
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t durationType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	_, err := time.ParseDuration(in.ValueString())
	if err != nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Invalid Duration String Value",
				"A string value was provided that is not valid duration.\n\n"+
					"Given Value: "+in.ValueString()+"\nError: "+err.Error(),
			),
		}
	}
	return durationValue{value: in.ValueString(), state: attr.ValueStateKnown}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t durationType) ValueFromTerraform(_ context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return DurationUnknown(), nil
	}
	if in.IsNull() {
		return DurationNull(), nil
	}

	var s string

	err := in.As(&s)
	if err != nil {
		return nil, err
	}

	_, err = time.ParseDuration(s)
	if err != nil {
		return nil, fmt.Errorf("unexpected error converting %s to %s: %w", in.Type(), t, err)
	}

	return durationValue{value: s, state: attr.ValueStateKnown}, nil
}

// Validate implements type validation. This type requires the value provided to be a StringValue value that is a parseable
// by time.Duration.
func (t durationType) Validate(_ context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Duration Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}

	if !in.IsKnown() || in.IsNull() {
		return
	}

	var strVal string

	if err := in.As(&strVal); err != nil {
		diags.AddAttributeError(
			path,
			"Duration Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}

	_, err := time.ParseDuration(strVal)
	if err != nil {
		diags.AddAttributeError(
			path,
			"Invalid Duration String Value",
			"A string value was provided that is not valid duration.\n\n"+
				"Given Value: "+strVal+"\n"+
				"Error: "+err.Error(),
		)
		return
	}

	return diags
}

func (t durationType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

func DurationValueFromDecimalSeconds(v decimal.Decimal) Duration {
	return DurationValue(time.Duration(int64(v.InexactFloat64() * float64(time.Second))))
}

func DurationValue(d time.Duration) Duration {
	return durationValue{
		value: d.String(),
		state: attr.ValueStateKnown,
	}
}

func DurationUnknown() Duration {
	return durationValue{
		state: attr.ValueStateUnknown,
	}
}

func DurationNull() Duration {
	return durationValue{
		state: attr.ValueStateNull,
	}
}

type Duration = durationValue

type durationValue struct {
	value string
	state attr.ValueState
}

func (d durationValue) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	if d.IsNull() {
		return tftypes.NewValue(tftypes.String, nil), nil
	}
	if d.IsUnknown() {
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue), nil
	}
	return tftypes.NewValue(tftypes.String, d.ValueString()), nil
}

func (d durationValue) Type(_ context.Context) attr.Type {
	return durationType{}
}

func (d durationValue) IsNull() bool {
	return d.state == attr.ValueStateNull
}

func (d durationValue) IsUnknown() bool {
	return d.state == attr.ValueStateUnknown
}

// String returns a human readable string of the value
func (d durationValue) String() string {
	var b strings.Builder
	b.WriteString(DurationType.String())
	b.WriteByte('<')
	switch d.state {
	case attr.ValueStateKnown:
		b.WriteString(d.ValueString())
	case attr.ValueStateNull:
		b.WriteString(attr.NullValueString)
	case attr.ValueStateUnknown:
		b.WriteString(attr.UnknownValueString)
	default:
		panic(fmt.Sprintf("unknown value state: %d", d.state))
	}
	b.WriteByte('>')
	return b.String()
}

// Equal returns true if the given value is equivalent.
func (d durationValue) Equal(x attr.Value) bool {
	o, ok := x.(durationValue)
	if !ok {
		return false
	}
	if d.IsNull() {
		return o.IsNull()
	}
	if d.IsUnknown() {
		return o.IsUnknown()
	}
	return d.value == o.value
}

// StringSemanticEquals returns true if the given string value can be parsed into a time.Duration which is equal to the durationValue.
func (d durationValue) StringSemanticEquals(ctx context.Context, in basetypes.StringValuable) (bool, diag.Diagnostics) {
	if d.IsNull() {
		return in.IsNull(), nil
	}
	if d.IsUnknown() {
		return in.IsUnknown(), nil
	}
	str, diags := in.ToStringValue(ctx)
	if diags.HasError() {
		return false, diags
	}
	dur, err := time.ParseDuration(str.ValueString())
	if err != nil {
		return false, nil
	}
	return d.ValueDuration() == dur, nil
}

func (d durationValue) ToStringValue(_ context.Context) (basetypes.StringValue, diag.Diagnostics) {
	if d.IsNull() {
		return basetypes.NewStringNull(), nil
	}
	if d.IsUnknown() {
		return basetypes.NewStringUnknown(), nil
	}
	return basetypes.NewStringValue(d.ValueString()), nil
}

func (d durationValue) ValueString() string {
	return d.value
}

func (d durationValue) ValueDuration() time.Duration {
	// it pretty much always should be parseable since the type implements validation
	v, _ := time.ParseDuration(d.value)
	return v
}
