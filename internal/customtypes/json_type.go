package customtypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*JsonType)(nil)
	_ xattr.TypeWithValidate  = (*JsonType)(nil)
)

// JsonType is an attribute type that represents a time duration.
type JsonType struct {
	basetypes.StringType
}

// String returns a human readable string of the type name.
func (t JsonType) String() string {
	return "extratypes.ScriptType"
}

// ValueType returns the Value type.
func (t JsonType) ValueType(_ context.Context) attr.Value {
	return Json{}
}

// Equal returns true if the given type is equivalent.
func (t JsonType) Equal(o attr.Type) bool {
	other, ok := o.(JsonType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// Validate implements type validation. This type requires the value provided to be a String value
// containing a parseable JSON.
func (t JsonType) Validate(_ context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Json Type Validation Error",
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
			"Json Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}

	if !json.Valid([]byte(strVal)) {
		diags.AddAttributeError(
			path,
			"Json Type Validation Error",
			"Value is not a valid json. This is an error in the configuration.",
		)
		return
	}

	return diags
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t JsonType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return Json{
		StringValue: in,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t JsonType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
