package customtypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*DurationType)(nil)
	_ xattr.TypeWithValidate  = (*DurationType)(nil)
)

// DurationType is an attribute type that represents a time duration.
type DurationType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t DurationType) String() string {
	return "extratypes.DurationType"
}

// ValueType returns the Value type.
func (t DurationType) ValueType(ctx context.Context) attr.Value {
	return Duration{}
}

// Equal returns true if the given type is equivalent.
func (t DurationType) Equal(o attr.Type) bool {
	other, ok := o.(DurationType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// Validate implements type validation. This type requires the value provided to be a String value that is a parseable
// by time.Duration.
func (t DurationType) Validate(ctx context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
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

// ValueFromString returns a StringValuable type given a StringValue.
func (t DurationType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return Duration{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t DurationType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
