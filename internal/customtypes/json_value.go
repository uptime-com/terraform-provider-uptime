package customtypes

import (
	"bytes"
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/nsf/jsondiff"
)

var (
	_ basetypes.StringValuableWithSemanticEquals = (*Json)(nil)
)

// Json represents a time duration.
type Json struct {
	basetypes.StringValue
}

// Type returns a JsonType.
func (v Json) Type(_ context.Context) attr.Type {
	return JsonType{}
}

// Equal returns true if the given type is equivalent.
func (v Json) Equal(other attr.Value) bool {
	o, ok := other.(Json)
	if !ok {
		return false
	}
	return v.StringValue.Equal(o.StringValue)
}

func (v Json) StringSemanticEquals(_ context.Context, val basetypes.StringValuable) (bool, diag.Diagnostics) {
	other, ok := val.(Json)
	if !ok {
		return false, nil
	}
	if v.IsNull() && other.IsNull() {
		return true, nil
	}
	if v.IsUnknown() || other.IsUnknown() {
		return true, nil
	}
	buf0 := bytes.NewBufferString(v.ValueString())
	buf1 := bytes.NewBufferString(other.ValueString())
	opts := jsondiff.DefaultJSONOptions()
	diff, _ := jsondiff.CompareStreams(buf0, buf1, &opts)
	return diff == jsondiff.FullMatch, nil
}

// NewJsonUnknown creates an Json with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewJsonUnknown() Json {
	return Json{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewJsonNull creates an Json with a null value. Determine whether the value is null via IsNull method.
func NewJsonNull() Json {
	return Json{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewJsonValue creates an Json with a known value.
func NewJsonValue(value string) Json {
	return Json{
		StringValue: basetypes.NewStringValue(value),
	}
}
