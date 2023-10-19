package reflect

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"reflect"

	"github.com/gobeam/stringy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/reflectwalk"
	"github.com/shopspring/decimal"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"
)

func CopyOut(dst any, src any) error {
	return reflectwalk.Walk(src, &pathWalker{
		PathWalker: &copyOutWalker{
			dst: dst,
		},
	})
}

type copyOutWalker struct {
	dst any
}

func (w *copyOutWalker) Walk(path string, tag Tag, f reflect.Value) error {
	// FIXME: this should more rely on TF typing system and less on native reflection
	if tag.Skip {
		return nil
	}
	if tag.Path != "" {
		path = tag.Path
	}
	a, ok := f.Interface().(attr.Value)
	if !ok {
		return nil // skip value which isn't an attribute
	}
	if a.IsUnknown() {
		return nil // skip value without value set
	}
	if a.IsNull() {
		return nil // skip null value
	}
	t, err := FindByPath(w.dst, path)
	if err != nil {
		if errors.Is(err, ErrNotFound) && tag.Opt {
			return nil
		}
		return err
	}
	if !t.CanSet() {
		return fmt.Errorf("value is not settable: %s", path)
	}
	err = w.copyOut(f, *t, tag)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

func (w *copyOutWalker) copyOut(f reflect.Value, t reflect.Value, tag Tag) (err error) {
	switch x := f.Interface().(type) {
	case types.String:
		if t.Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T", t.Interface())
		}
		t.SetString(x.ValueString())
	case types.Bool:
		if t.Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T", t.Interface())
		}
		t.SetBool(x.ValueBool())
	case types.Int64:
		if t.Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T", t.Interface())
		}
		t.SetInt(x.ValueInt64())
	case types.Float64:
		if t.Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T", t.Interface())
		}
		t.SetFloat(x.ValueFloat64())
	case types.List:
		if t.Kind() != reflect.Slice {
			return fmt.Errorf("slice expected, got %T", t.Interface())
		}
		err = w.sliceTo(x, t)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case types.Set:
		if t.Kind() != reflect.Slice {
			return fmt.Errorf("slice expected, got %T", t.Interface())
		}
		err = w.sliceTo(x, t)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case types.Map:
		switch tag.Extra {
		case "headers":
			if t.Kind() != reflect.String {
				return fmt.Errorf("string expected, got %T", t.Interface())
			}
			err = w.copyOutHeadersMap(x, t)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		default:
			if t.Kind() != reflect.Map {
				return fmt.Errorf("map expected, got %T", t.Interface())
			}
			err = w.copyOutMap(x, t)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	case types.Object:
		if t.Kind() != reflect.Struct {
			return fmt.Errorf("struct expected, got %T", t.Interface())
		}
		err = w.objectTo(x, t)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case types.Number:
		return w.copyOutNumber(x, t)
	case customtypes.Duration:
		return w.copyOutDuration(x, t)
	case customtypes.SmartPercentage:
		return w.copyOutSmartPercentage(x, t)
	case customtypes.Json:
		if t.Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T", t.Interface())
		}
		t.SetString(x.ValueString())
	default:
		return fmt.Errorf("unsupported type: %T", x)
	}
	return nil
}

type sliceElementable interface {
	Elements() []attr.Value
	ElementType(ctx context.Context) attr.Type
}

func (w *copyOutWalker) sliceTo(x sliceElementable, v reflect.Value) error {
	els := x.Elements()
	v.Grow(len(els))
	v.SetLen(len(els))
	switch et := x.ElementType(nil); et {
	case types.StringType:
		if v.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T", v.Interface())
		}
		for i, e := range x.Elements() {
			v.Index(i).Set(reflect.ValueOf(e.(types.String).ValueString()))
		}
	case types.BoolType:
		if v.Type().Elem().Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T", v.Interface())
		}
	case types.Int64Type:
		if v.Type().Elem().Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T", v.Interface())
		}
	case types.Float64Type:
		if v.Type().Elem().Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T", v.Interface())
		}
	default:
		return fmt.Errorf("unsupported slice or set element type: %T", et)
	}
	return nil
}

type mapElementable interface {
	Elements() map[string]attr.Value
	ElementType(ctx context.Context) attr.Type
}

func (w *copyOutWalker) copyOutMap(x mapElementable, v reflect.Value) error {
	els := x.Elements()
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	switch et := x.ElementType(nil); et {
	case types.StringType:
		if v.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T", v.Interface())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.String).ValueString()))
		}
	case types.BoolType:
		if v.Type().Elem().Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T", v.Interface())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.Bool).ValueBool()))
		}
	case types.Int64Type:
		if v.Type().Elem().Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T", v.Interface())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.Int64).ValueInt64()))
		}
	case types.Float64Type:
		if v.Type().Elem().Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T", v.Interface())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.Float64).ValueFloat64()))
		}
	default:
		return fmt.Errorf("unsupported map element type: %T", et)
	}
	return nil
}

func (w *copyOutWalker) copyOutHeadersMap(x mapElementable, v reflect.Value) error {
	header := make(http.Header)
	for key, val := range x.Elements() {
		values, ok := val.(types.List)
		if !ok {
			return fmt.Errorf("types.ListType{types.StringType} element expected, got %T", val)
		}
		if values.ElementType(nil) != types.StringType {
			return fmt.Errorf("types.ListType{types.StringType} element expected, got %T", val)
		}
		for _, e := range values.Elements() {
			header.Add(key, e.(types.String).ValueString())
		}
	}
	buf := bytes.NewBuffer(nil)
	if err := header.Write(buf); err != nil {
		return nil
	}
	v.SetString(buf.String())
	return nil
}

type objectAttributable interface {
	Attributes() map[string]attr.Value
}

func (w *copyOutWalker) objectTo(f objectAttributable, t reflect.Value) error {
	for attrName, attrVal := range f.Attributes() {
		fn := stringy.New(attrName).CamelCase(
			"id", "ID",
			"url", "URL",
			"crl", "CRL",
		)
		fv := t.FieldByName(fn)
		if !fv.IsValid() {
			return fmt.Errorf("field %s is not found", fn)
		}
		if !fv.CanSet() {
			return fmt.Errorf("field %s is not settable", fn)
		}
		switch x := attrVal.(type) {
		case types.String:
			if fv.Type().Kind() != reflect.String {
				return fmt.Errorf("string expected, got %T", fv.Interface())
			}
			fv.SetString(x.ValueString())
		case types.Bool:
			if fv.Type().Kind() != reflect.Bool {
				return fmt.Errorf("bool expected, got %T", fv.Interface())
			}
			fv.SetBool(x.ValueBool())
		case types.Int64:
			if fv.Type().Kind() != reflect.Int {
				return fmt.Errorf("int expected, got %T", fv.Interface())
			}
			fv.SetInt(x.ValueInt64())
		case types.Float64:
			if fv.Type().Kind() != reflect.Float64 {
				return fmt.Errorf("float64 expected, got %T", fv.Interface())
			}
			fv.SetFloat(x.ValueFloat64())
		case types.List, types.Set:
			if fv.Type().Kind() != reflect.Slice {
				return fmt.Errorf("slice expected, got %T", fv.Interface())
			}
			if err := w.sliceTo(x.(sliceElementable), fv); err != nil {
				return err
			}
		case types.Map:
			if fv.Type().Kind() != reflect.Map {
				return fmt.Errorf("map expected, got %T", fv.Interface())
			}
			if err := w.copyOutMap(x, fv); err != nil {
				return err
			}
		case types.Object:
			if fv.Type().Kind() != reflect.Struct {
				return fmt.Errorf("struct expected, got %T", fv.Interface())
			}
			if err := w.objectTo(x, fv); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %T", x)
		}
	}
	return nil
}

func (w *copyOutWalker) derefInit(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return w.derefInit(v.Elem())
	}
	return v
}

func (w *copyOutWalker) copyOutNumber(x types.Number, val reflect.Value) error {
	val = w.derefInit(val)
	val.Set(reflect.ValueOf(decimal.RequireFromString(x.ValueBigFloat().String())))
	return nil
}

func (w *copyOutWalker) copyOutDuration(x customtypes.Duration, val reflect.Value) error {
	val = w.derefInit(val)
	dur, diags := x.ValueDuration()
	if diags.HasError() {
		return fmt.Errorf("%ss: %s", diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
	}
	val.Set(reflect.ValueOf(decimal.NewFromFloat(dur.Seconds())))
	return nil
}

func (w *copyOutWalker) copyOutSmartPercentage(x customtypes.SmartPercentage, val reflect.Value) error {
	i, _ := new(big.Float).Mul(x.ValueBigFloat(), big.NewFloat(10000)).Int(nil)
	d := decimal.NewFromBigInt(i, -4)
	w.derefInit(val).Set(reflect.ValueOf(d))
	return nil
}
