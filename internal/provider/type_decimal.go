package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/shopspring/decimal"
)

var DecimalType = decimalType{}

type decimalType struct{}

func (t decimalType) String() string {
	return "provider.DecimalType"
}

func (t decimalType) TerraformType(context.Context) tftypes.Type {
	return tftypes.String
}

func (t decimalType) ValueType(context.Context) attr.Value {
	return decimalValue{}
}

func (t decimalType) Equal(other attr.Type) bool {
	_, ok := other.(decimalType)
	return ok
}

func (t decimalType) ValueFromString(_ context.Context, v basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	d, err := decimal.NewFromString(v.ValueString())
	if err != nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic("Invalid Decimal Value",
				fmt.Sprintf("A string value was provided that is not valid decimal.\n\nGiven Value: %s\nError: %s", v.ValueString(), err.Error()),
			),
		}
	}
	return DecimalValue(d), nil
}

func (t decimalType) ValueFromTerraform(_ context.Context, value tftypes.Value) (attr.Value, error) {
	if value.IsNull() {
		return decimalValue{state: attr.ValueStateNull}, nil
	}
	if !value.IsKnown() {
		return decimalValue{state: attr.ValueStateUnknown}, nil
	}
	var str string
	err := value.As(&str)
	if err != nil {
		return nil, err
	}
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return nil, err
	}
	return decimalValue{
		value: dec,
		state: attr.ValueStateKnown,
	}, nil
}

func (t decimalType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

func DecimalValue(v decimal.Decimal) Decimal {
	return decimalValue{
		value: v,
		state: attr.ValueStateKnown,
	}
}

func DecimalNull() Decimal {
	return decimalValue{state: attr.ValueStateNull}
}

func DecimalUnknown() Decimal {
	return decimalValue{state: attr.ValueStateUnknown}
}

type Decimal = decimalValue

type decimalValue struct {
	value decimal.Decimal
	state attr.ValueState
}

func (d decimalValue) Type(context.Context) attr.Type {
	return DecimalType
}

func (d decimalValue) IsNull() bool {
	return d.state == attr.ValueStateNull
}

func (d decimalValue) IsUnknown() bool {
	return d.state == attr.ValueStateUnknown
}

func (d decimalValue) String() string {
	if d.IsNull() {
		return attr.NullValueString
	}
	if d.IsUnknown() {
		return attr.UnknownValueString
	}
	return d.value.String()
}

func (d decimalValue) ValueDecimal() decimal.Decimal {
	return d.value
}

func (d decimalValue) ToTerraformValue(context.Context) (tftypes.Value, error) {
	return tftypes.NewValue(tftypes.String, d.value.String()), nil
}

func (d decimalValue) Equal(x attr.Value) bool {
	o, ok := x.(decimalValue)
	if !ok {
		return false
	}
	if d.IsNull() {
		return o.IsNull()
	}
	if d.IsUnknown() {
		return o.IsUnknown()
	}
	return d.value.Equal(o.value)
}

func (d decimalValue) StringSemanticEquals(ctx context.Context, v basetypes.StringValuable) (bool, diag.Diagnostics) {
	sv, diags := v.ToStringValue(ctx)
	if diags.HasError() {
		return false, diags
	}
	dv, err := decimal.NewFromString(sv.ValueString())
	if err != nil {
		return false, nil
	}
	return d.value.Equal(dv), nil
}

func (d decimalValue) ToStringValue(context.Context) (basetypes.StringValue, diag.Diagnostics) {
	return basetypes.NewStringValue(d.value.String()), nil
}
