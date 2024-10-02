package provider

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SetPrimitives interface {
	string | int64 | float64 | int32 | bool
}

type SetAttributeAdapter[P SetPrimitives] struct{}

func (a SetAttributeAdapter[P]) Slice(v types.Set) []P {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	elems := v.Elements()
	res := make([]P, len(elems))
	for i := range elems {
		switch elems[i].Type(nil) {
		case types.StringType:
			res[i] = any(elems[i].(types.String).ValueString()).(P)
		case types.BoolType:
			res[i] = any(elems[i].(types.Bool).ValueBool()).(P)
		case types.Int32Type:
			res[i] = any(elems[i].(types.Int32).ValueInt32()).(P)
		case types.Int64Type:
			res[i] = any(elems[i].(types.Int64).ValueInt64()).(P)
		case types.Float64Type:
			res[i] = any(elems[i].(types.Float64).ValueFloat64()).(P)
		default:
			panic("unsupported set element type")
		}
	}
	return res
}

func (a SetAttributeAdapter[P]) SliceValue(p []P) types.Set {
	var elems []attr.Value
	for i := range p {
		switch x := any(p[i]).(type) {
		case string:
			elems = append(elems, types.StringValue(x))
		case bool:
			elems = append(elems, types.BoolValue(x))
		case int32:
			elems = append(elems, types.Int32Value(x))
		case int64:
			elems = append(elems, types.Int64Value(x))
		case float64:
			elems = append(elems, types.Float64Value(x))
		default:
			panic("unsupported set element type")
		}
	}
	switch any(p).(type) {
	case []string:
		return types.SetValueMust(types.StringType, elems)
	case []bool:
		return types.SetValueMust(types.BoolType, elems)
	case []int32:
		return types.SetValueMust(types.Int32Type, elems)
	case []int64:
		return types.SetValueMust(types.Int64Type, elems)
	case []float64:
		return types.SetValueMust(types.Float64Type, elems)
	default:
		panic("unsupported set element type")
	}
}

type zoyaDescriber struct{}

func (e zoyaDescriber) Description(context.Context) string {
	return ""
}

func (e zoyaDescriber) MarkdownDescription(context.Context) string {
	return ""
}

func ErrorAccumulator[T any](acc *multierror.Error) func(T, error) T {
	return func(v T, err error) T {
		if err != nil {
			acc = multierror.Append(acc, err)
		}
		return v
	}
}
