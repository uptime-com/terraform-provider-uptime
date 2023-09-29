package customtypes

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.NumberTypable = (*SmartPercentageType)(nil)
)

// SmartPercentageType represents a percentage in form of floating point number value with the following semantic and
// validation sugar:
//
// - Valid range for percentage is [0, 100], negative values are not allowed, values over 100 are not allowed;
// - Input values in range [0, 1) are interpreted as fractions of 100, e.g. 0.5 is interpreted as 50%;
// - Input values in range [1, 100] are interpreted as percentages, e.g. 1 is interpreted as 1%, 50 is interpreted as 50%;
// - For fractional representation, values with more than 4 digits after decimal point are not allowed;
// - For percentage representation, values with more than 2 digits after decimal point are not allowed;
type SmartPercentageType struct {
	basetypes.NumberType
}

// String returns a human readable string of the type name
func (t SmartPercentageType) String() string {
	return "customtypes.SmartPercentageType"
}

// ValueType returns the Value type.
func (t SmartPercentageType) ValueType(_ context.Context) attr.Value {
	return SmartPercentage{}
}

// Equal returns true if the given type is equivalent.
func (t SmartPercentageType) Equal(o attr.Type) bool {
	other, ok := o.(SmartPercentageType)

	if !ok {
		return false
	}

	return t.NumberType.Equal(other.NumberType)
}

// Validate implements type validation
func (t SmartPercentageType) Validate(_ context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return diags
	}
	if !in.Type().Equal(tftypes.Number) {
		diags.AddAttributeError(
			path,
			"SmartPercentage Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Expected Number value, received %T with value: %v", in, in),
		)
		return diags
	}
	if !in.IsKnown() || in.IsNull() {
		return diags
	}

	var value *big.Float
	err := in.As(&value)

	if err != nil {
		diags.AddAttributeError(
			path,
			"SmartPercentage Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Cannot convert value to big.Float: %s", err),
		)
		return diags
	}

	valueFloat64, _ := value.Float64()
	if valueFloat64 < 0 || valueFloat64 > 100 {
		diags.AddAttributeError(
			path,
			"SmartPercentage Type Validation Error",
			"Value must be in range [0, 100]. This is an error in the configuration.",
		)
		return diags
	}
	if valueFloat64 < 1 {
		if math.Round(valueFloat64*10000) != valueFloat64*10000 {
			diags.AddAttributeError(
				path,
				"SmartPercentage Type Validation Error",
				"Value must have at most 4 digits after decimal point if presented as fraction. This is an error in the configuration.",
			)
			return
		}
	} else {
		if math.Round(valueFloat64*100) != valueFloat64*100 {
			diags.AddAttributeError(
				path,
				"SmartPercentage Type Validation Error",
				"Value must have at most 2 digits after decimal point if presented as percentage. This is an error in the configuration.",
			)
			return diags
		}
	}

	return
}

func (t SmartPercentageType) ValueFromFloat64(_ context.Context, in basetypes.NumberValue) (basetypes.NumberValuable, diag.Diagnostics) {
	return SmartPercentage{
		NumberValue: in,
	}, nil
}

func (t SmartPercentageType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return NewSmartPercentageUnknown(), nil
	}
	if in.IsNull() {
		return NewSmartPercentageNull(), nil
	}

	val, err := t.NumberType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	num := val.(basetypes.NumberValue)
	return NewSmartPercentageValue(num.ValueBigFloat()), nil
}
