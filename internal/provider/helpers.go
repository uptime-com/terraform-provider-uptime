package provider

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ptr[T any](v T) *T {
	return &v
}

func mustErr[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func mustDiag[T any](x T, diags diag.Diagnostics) T {
	if diags.HasError() {
		panic(diags)
	}
	return x
}

// mirror copies the fields of src into dst. dst must be a pointer to a struct. At least one of the structs must be
// defined in the current package. Local struct fields are matched with the same name in the other struct. That can
// be overridden with the "api" tag. Special case is a field tagged with api:"Pk". In addition to mapping it is also
// converted *string <-> *int.
func mirror(dst any, src any) {
	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("expected pointer to struct, got %T", dst))
	}
	for dv.Kind() == reflect.Ptr {
		dv = dv.Elem()
	}
	sv := reflect.ValueOf(src)
	for sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	// reference type is the type defined in the current package
	var ref reflect.Type
	switch {
	case strings.HasSuffix(dv.Type().PkgPath(), "provider"):
		ref = dv.Type()
	case strings.HasSuffix(sv.Type().PkgPath(), "provider"):
		ref = sv.Type()
	default:
		panic(fmt.Sprintf("expected at least one struct from provider package, got %T and %T", dst, src))
	}
	for i := 0; i < ref.NumField(); i++ {
		name := ref.Field(i).Name
		apiName, ok := ref.Field(i).Tag.Lookup("api")
		if !ok {
			apiName = name
		}
		var (
			dstName string
			srcName string
		)
		if ref == dv.Type() { // dst is of local type
			dstName = name
			srcName = apiName
		} else {
			dstName = apiName
			srcName = name
		}
		if sv.FieldByName(srcName).IsZero() {
			continue // skip zero values
		}
		if apiName == "Pk" {
			if ref == dv.Type() { // dst is of local type
				id := strconv.Itoa(sv.FieldByName(srcName).Elem().Interface().(int))
				dv.FieldByName(dstName).Set(reflect.ValueOf(ptr(id)))
			} else {
				pk, err := strconv.Atoi(sv.FieldByName(srcName).Elem().Interface().(string))
				if err != nil {
					panic(err)
				}
				dv.FieldByName(dstName).Set(reflect.ValueOf(ptr(pk)))
			}
		} else {
			dv.FieldByName(dstName).Set(sv.FieldByName(srcName))
		}
	}
}

func fromTerraform(ctx context.Context, dst any, src any) (d diag.Diagnostics) {
	dv := reflect.ValueOf(dst)
	for dv.Kind() == reflect.Ptr {
		dv = dv.Elem()
	}
	sv := reflect.ValueOf(src)
	for sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	st := sv.Type()

	setFunc := func(v reflect.Value, x reflect.Value) {
		if v.Kind() == reflect.Ptr {
			v.Set(reflect.New(v.Type().Elem()))
			v.Elem().Set(x)
		} else {
			v.Set(x)
		}
	}

	for i := 0; i < st.NumField(); i++ {
		if st.Field(i).Anonymous {
			continue
		}
		tag, ok := st.Field(i).Tag.Lookup("map")
		if !ok {
			continue // skip fields without map tag
		}
		dvf := dv.FieldByName(tag)
		tfv, ok := sv.FieldByName(st.Field(i).Name).Interface().(attr.Value)
		if !ok {
			continue // skip non-Value fields
		}
		if tfv.IsUnknown() {
			continue // skip unknown values
		}
		if tfv.IsNull() {
			if dvf.Kind() == reflect.Ptr {
				continue // skip null values for pointer fields...
			}
			d.Append(diag.NewErrorDiagnostic(
				"cannot set null value to non-pointer field",
				fmt.Sprintf("cannot set null value to non-pointer field %s", tag),
			))
			continue // ...or emit error for non-pointer fields
		}
		switch tfv.Type(ctx) {
		case types.BoolType:
			v, dd := tfv.(types.Bool).ToBoolValue(ctx)
			if dd.HasError() {
				d.Append(dd...)
				continue
			}
			setFunc(dvf, reflect.ValueOf(v.ValueBool()))
		case types.StringType:
			v, dd := tfv.(types.String).ToStringValue(ctx)
			if dd.HasError() {
				d.Append(dd...)
				continue
			}
			setFunc(dvf, reflect.ValueOf(v.ValueString()))
		case types.Int64Type:
			v, dd := tfv.(types.Int64).ToInt64Value(ctx)
			if dd.HasError() {
				d.Append(dd...)
				continue
			}
			setFunc(dvf, reflect.ValueOf(v.ValueInt64()))
		case types.NumberType:
			v, dd := tfv.(types.Number).ToNumberValue(ctx)
			if dd.HasError() {
				d.Append(dd...)
				continue
			}
			f64, _ := v.ValueBigFloat().Float64()
			setFunc(dvf, reflect.ValueOf(f64))
		case types.SetType{ElemType: types.StringType}:
			v, dd := tfv.(types.Set).ToSetValue(ctx)
			if dd.HasError() {
				d.Append(dd...)
				continue
			}
			x := make([]string, 0)
			dd = v.ElementsAs(ctx, &x, false)
			if dd.HasError() {
				d.Append(dd...)
				continue
			}
			setFunc(dvf, reflect.ValueOf(x))
		default:
			panic(fmt.Sprintf("unsupported Terraform type %T", tfv))
		}
	}
	return
}

func toTerraform(dst any, src any) (d diag.Diagnostics) {
	dv := reflect.ValueOf(dst)
	for dv.Kind() == reflect.Ptr {
		dv = dv.Elem()
	}
	dt := dv.Type()
	sv := reflect.ValueOf(src)
	for sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	for i := 0; i < dt.NumField(); i++ {
		if dt.Field(i).Anonymous {
			continue
		}
		tag, ok := dt.Field(i).Tag.Lookup("map")
		if !ok {
			continue // skip fields without map tag
		}
		sf := sv.FieldByName(tag)
		if sf.Kind() == reflect.Ptr && sf.IsZero() {
			continue // nil pointers leave target fields unset
		}
		for sf.Kind() == reflect.Ptr {
			sf = sf.Elem()
		}
		df := dv.FieldByName(dt.Field(i).Name)
		switch x := df.Interface().(type) {
		case basetypes.StringValue:
			switch y := sf.Interface().(type) {
			case string:
				x = basetypes.NewStringValue(y)
			case *string:
				if y != nil {
					x = basetypes.NewStringValue(*y)
				} else {
					x = basetypes.NewStringNull()
				}
			default:
				d.AddError("unsupported type mapping", fmt.Sprintf("unsupported type mapping tf:%T -> %T", x, y))
				return
			}
			df.Set(reflect.ValueOf(x))
		case basetypes.Int64Value:
			switch y := sf.Interface().(type) {
			case int:
				x = basetypes.NewInt64Value(int64(y))
			case *int:
				if y != nil {
					x = basetypes.NewInt64Value(int64(*y))
				} else {
					x = basetypes.NewInt64Null()
				}
			case int64:
				x = basetypes.NewInt64Value(y)
			case *int64:
				if y != nil {
					x = basetypes.NewInt64Value(*y)
				} else {
					x = basetypes.NewInt64Null()
				}
			default:
				d.AddError("unsupported type mapping", fmt.Sprintf("unsupported type mapping tf:%T -> %T", x, y))
				return
			}
			df.Set(reflect.ValueOf(x))
		case basetypes.BoolValue:
			switch y := sf.Interface().(type) {
			case bool:
				x = basetypes.NewBoolValue(y)
			case *bool:
				if y != nil {
					x = basetypes.NewBoolValue(*y)
				} else {
					x = basetypes.NewBoolNull()
				}
			default:
				d.AddError("unsupported type mapping", fmt.Sprintf("unsupported type mapping tf:%T -> %T", x, y))
				return
			}
			df.Set(reflect.ValueOf(x))
		case basetypes.SetValue:
			a := make([]string, 0)
			dd := diag.Diagnostics{}
			switch y := sf.Interface().(type) {
			case []string:
				a = y
			case *[]string:
				a = *y
			}
			e := make([]attr.Value, len(a))
			for k := range a {
				e[k] = basetypes.NewStringValue(a[k])
			}
			x, dd = basetypes.NewSetValue(basetypes.StringType{}, e)
			d.Append(dd...)
			if d.HasError() {
				return
			}
			df.Set(reflect.ValueOf(x))
		default:
			d.AddError("unsupported type", fmt.Sprintf("unsupported type: %T", x))
			return
		}
	}
	return
}
