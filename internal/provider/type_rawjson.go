package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/nsf/jsondiff"
)

var (
	_ xattr.TypeWithValidate  = (*RawJsonType)(nil)
	_ basetypes.StringTypable = (*RawJsonType)(nil)
)

// RawJsonType is an attribute type that represents a time duration.
type RawJsonType struct{}

// StringValue returns a human readable string of the type name.
func (t RawJsonType) String() string {
	return "RawJsonType"
}

// ValueType returns the Value type.
func (t RawJsonType) ValueType(_ context.Context) attr.Value {
	return rawJsonValue{}
}

// Equal returns true if the given type is equivalent.
func (t RawJsonType) Equal(o attr.Type) bool {
	_, ok := o.(RawJsonType)
	return ok
}

// Validate implements type validation. This type requires the value provided to be a StringValue value
// containing a parseable JSON.
func (t RawJsonType) Validate(_ context.Context, in tftypes.Value, path path.Path) (diags diag.Diagnostics) {
	if in.Type() == nil {
		return
	}

	if !in.Type().Is(tftypes.String) {
		err := fmt.Errorf("expected String value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			fmt.Sprintf("%s Validation Error", t.String()),
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
			fmt.Sprintf("%s Validation Error", t.String()),
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return
	}

	if !json.Valid([]byte(strVal)) {
		diags.AddAttributeError(
			path,
			fmt.Sprintf("%s Validation Error", t.String()),
			"Value is not a valid json. This is an error in the configuration.",
		)
		return
	}

	return diags
}

func (t RawJsonType) TerraformType(context.Context) tftypes.Type {
	return tftypes.String
}

func (t RawJsonType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t RawJsonType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return RawJsonValue(in.ValueString()), nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t RawJsonType) ValueFromTerraform(_ context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return RawJsonUnknown(), nil
	}
	if in.IsNull() {
		return RawJsonNull(), nil
	}
	var strVal string
	if err := in.As(&strVal); err != nil {
		return nil, err
	}
	return RawJsonValue(strVal), nil
}

// RawJsonNull creates an rawJsonValue with a null value. Determine whether the value is null via IsNull method.
func RawJsonNull() RawJson {
	return rawJsonValue{state: attr.ValueStateNull}
}

// RawJsonUnknown creates an rawJsonValue with an unknown value. Determine whether the value is unknown via IsUnknown method.
func RawJsonUnknown() RawJson {
	return rawJsonValue{state: attr.ValueStateUnknown}
}

// RawJsonValue creates an rawJsonValue with a known value.
func RawJsonValue(value string) RawJson {
	return rawJsonValue{state: attr.ValueStateKnown, value: value}
}

var (
	_ attr.Value                                 = (*rawJsonValue)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*rawJsonValue)(nil)
)

type RawJson = rawJsonValue

// rawJsonValue represents a time duration.
type rawJsonValue struct {
	state attr.ValueState
	value string
}

// String returns a human readable string of the value.
func (v rawJsonValue) String() string {
	return "moretypes.RawJsonValue<...>" // TODO: implement better representation
}

// ValueString returns the value as a string.
func (v rawJsonValue) ValueString() string {
	return v.value
}

// Type returns a RawJsonType.
func (v rawJsonValue) Type(context.Context) attr.Type {
	return RawJsonType{}
}

// IsNull returns true if the value is null.
func (v rawJsonValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

// IsUnknown returns true if the value is unknown.
func (v rawJsonValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

// Equal returns true if the given type is equivalent.
func (v rawJsonValue) Equal(other attr.Value) bool {
	o, ok := other.(rawJsonValue)
	if !ok {
		return false
	}
	return v.value == o.value
}

// StringSemanticEquals returns true if JSON values are equal disregarding formatting.
func (v rawJsonValue) StringSemanticEquals(_ context.Context, val basetypes.StringValuable) (bool, diag.Diagnostics) {
	o, ok := val.(rawJsonValue)
	if !ok {
		return false, nil
	}
	if v.IsNull() && o.IsNull() {
		return true, nil
	}
	if v.IsUnknown() && o.IsUnknown() {
		return true, nil
	}
	buf0 := bytes.NewBufferString(v.ValueString())
	buf1 := bytes.NewBufferString(o.ValueString())
	opts := jsondiff.DefaultJSONOptions()
	diff, _ := jsondiff.CompareStreams(buf0, buf1, &opts)
	return diff == jsondiff.FullMatch, nil
}

// ToTerraformValue returns a tftypes.Value representation of the value.
func (v rawJsonValue) ToTerraformValue(context.Context) (tftypes.Value, error) {
	return tftypes.NewValue(tftypes.String, v.value), nil
}

// ToStringValue returns a basetypes.StringValue representation of the value.
func (v rawJsonValue) ToStringValue(context.Context) (basetypes.StringValue, diag.Diagnostics) {
	return basetypes.NewStringValue(v.value), nil
}
