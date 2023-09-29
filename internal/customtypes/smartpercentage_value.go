package customtypes

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ basetypes.NumberValuable                   = (*SmartPercentage)(nil)
	_ basetypes.NumberValuableWithSemanticEquals = (*SmartPercentage)(nil)
)

// SmartPercentage represents a value of SmartPercentageType.
type SmartPercentage struct {
	basetypes.NumberValue
}

// Type returns a SmartPercentageType.
func (v SmartPercentage) Type(_ context.Context) attr.Type {
	return SmartPercentageType{}
}

func (v SmartPercentage) Equal(value attr.Value) bool {
	other, ok := value.(SmartPercentage)
	if !ok {
		return false
	}
	return v.NumberValue.Equal(other.NumberValue)
}

func (v SmartPercentage) NumberSemanticEquals(ctx context.Context, val basetypes.NumberValuable) (res bool, diag diag.Diagnostics) {
	other, ok := val.(SmartPercentage)
	if !ok {
		return false, nil
	}
	if v.IsNull() && other.IsNull() {
		return true, nil
	}
	if v.IsUnknown() || other.IsUnknown() {
		return true, nil
	}
	num0, num0Diag := v.ToNumberValue(ctx)
	diag.Append(num0Diag...)

	num1, num1Diag := other.ToNumberValue(ctx)
	diag.Append(num1Diag...)

	if diag.HasError() {
		return false, diag
	}

	big0, _ := new(big.Float).Mul(num0.ValueBigFloat(), big.NewFloat(10000)).Int(nil)
	big1, _ := new(big.Float).Mul(num1.ValueBigFloat(), big.NewFloat(10000)).Int(nil)

	res = big0.Cmp(big1) == 0
	return res, nil
}

func (v SmartPercentage) ValueBigFloat() *big.Float {
	if v.IsNull() || v.IsUnknown() {
		return big.NewFloat(0)
	}
	return v.NumberValue.ValueBigFloat()
}

// NewSmartPercentageUnknown creates an SmartPercentage with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewSmartPercentageUnknown() SmartPercentage {
	return SmartPercentage{
		NumberValue: basetypes.NewNumberUnknown(),
	}
}

// NewSmartPercentageNull creates an SmartPercentage with a null value. Determine whether the value is null via IsNull method.
func NewSmartPercentageNull() SmartPercentage {
	return SmartPercentage{
		NumberValue: basetypes.NewNumberNull(),
	}
}

// NewSmartPercentageValue creates an SmartPercentage with a known value.
func NewSmartPercentageValue(value *big.Float) SmartPercentage {
	if value.Cmp(big.NewFloat(1)) >= 0 {
		return NewSmartPercentageValuePercentage(value)
	}
	return NewSmartPercentageValueFraction(value)
}

func NewSmartPercentageValuePercentage(value *big.Float) SmartPercentage {
	return SmartPercentage{
		NumberValue: basetypes.NewNumberValue(value.Quo(value, big.NewFloat(100))),
	}
}

func NewSmartPercentageValueFraction(value *big.Float) SmartPercentage {
	return SmartPercentage{
		NumberValue: basetypes.NewNumberValue(value),
	}
}
