package reflect

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"reflect"
	"strings"
	"time"

	"github.com/gobeam/stringy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/reflectwalk"
	"github.com/shopspring/decimal"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"
)

func CopyIn(dst any, src any) error {
	return reflectwalk.Walk(dst, &pathWalker{
		PathWalker: &copyInWalker{
			src: src,
		},
	})
}

type copyInWalker struct {
	src any
}

func (w *copyInWalker) Walk(path string, tag Tag, t reflect.Value) error {
	// FIXME: this should more rely on TF typing system and less on native reflection
	if _, ok := t.Interface().(attr.Value); !ok {
		return nil // skip value which isn't an attribute
	}
	if tag.Skip {
		return nil
	}
	if tag.Path != "" {
		path = tag.Path
	}
	f, err := FindByPath(w.src, path)
	if err != nil {
		if errors.Is(err, ErrNotFound) && tag.Opt {
			return nil
		}
		return err
	}
	err = w.copyIn(*f, t, tag)
	if err != nil {
		return fmt.Errorf("%s: copy in error: %w", path, err)
	}
	return nil
}

func (w *copyInWalker) copyIn(f reflect.Value, t reflect.Value, tag Tag) error {
	switch t.Interface().(type) {
	case types.String:
		if f.Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T", f.Interface())
		}
		t.Set(reflect.ValueOf(types.StringValue(f.String())))
	case types.Bool:
		if f.Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T", f.Interface())
		}
		t.Set(reflect.ValueOf(types.BoolValue(f.Bool())))
	case types.Int64:
		if f.Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T", f.Interface())
		}
		t.Set(reflect.ValueOf(types.Int64Value(f.Int())))
	case types.Float64:
		if f.Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T", f.Interface())
		}
		t.Set(reflect.ValueOf(types.Float64Value(f.Float())))
	case types.List, types.Set:
		err := w.copyInSlice(f, t)
		if err != nil {
			return err
		}
	case types.Map:
		switch tag.Extra {
		case "headers":
			err := w.copyInHeadersMap(f, t)
			if err != nil {
				return err
			}
		default:
			return w.copyInMap(f, t)
		}
	case types.Number:
		err := w.copyInNumber(f, t)
		if err != nil {
			return err
		}
	case types.Object:
		err := w.copyInObject(f, t)
		if err != nil {
			return err
		}
	case customtypes.Duration:
		err := w.copyInDuration(f, t)
		if err != nil {
			return err
		}
	case customtypes.SmartPercentage:
		err := w.copyInSmartPercentage(f, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *copyInWalker) copyInSlice(f, t reflect.Value) error {
	if f.Kind() != reflect.Slice {
		return fmt.Errorf("slice expected, got %T", f.Interface())
	}
	var (
		typ attr.Type
		els []attr.Value
	)
	switch f.Type().Elem().Kind() {
	case reflect.String:
		typ = types.StringType
		for i := 0; i < f.Len(); i++ {
			els = append(els, types.StringValue(f.Index(i).String()))
		}
	case reflect.Bool:
		typ = types.BoolType
		for i := 0; i < f.Len(); i++ {
			els = append(els, types.BoolValue(f.Index(i).Bool()))
		}
	case reflect.Int:
		typ = types.Int64Type
		for i := 0; i < f.Len(); i++ {
			els = append(els, types.Int64Value(f.Index(i).Int()))
		}
	case reflect.Float64:
		typ = types.Float64Type
		for i := 0; i < f.Len(); i++ {
			els = append(els, types.Float64Value(f.Index(i).Float()))
		}
	default:
		return fmt.Errorf("unsupported slice element type: %T", f.Interface())
	}
	switch t.Interface().(type) {
	case types.List:
		t.Set(reflect.ValueOf(types.ListValueMust(typ, els)))
	case types.Set:
		t.Set(reflect.ValueOf(types.SetValueMust(typ, els)))
	default:
		return fmt.Errorf("unsupported destination type: %T", t.Interface())
	}
	return nil
}

func (w *copyInWalker) copyInMap(f, t reflect.Value) error {
	if f.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("unsupported map key type: %T", f.Interface())
	}
	var (
		typ attr.Type
		els = make(map[string]attr.Value)
	)
	switch f.Type().Elem().Kind() {
	case reflect.String:
		typ = types.StringType
		for _, k := range f.MapKeys() {
			els[k.String()] = types.StringValue(f.MapIndex(k).String())
		}
	case reflect.Bool:
		typ = types.BoolType
		for _, k := range f.MapKeys() {
			els[f.String()] = types.BoolValue(f.MapIndex(k).Bool())
		}
	case reflect.Int:
		typ = types.Int64Type
		for _, k := range f.MapKeys() {
			els[k.String()] = types.Int64Value(f.MapIndex(k).Int())
		}
	case reflect.Float64:
		typ = types.Float64Type
		for _, k := range f.MapKeys() {
			els[k.String()] = types.Float64Value(f.MapIndex(k).Float())
		}
	default:
		return fmt.Errorf("unsupported map element type: %T", f.Interface())
	}
	t.Set(reflect.ValueOf(types.MapValueMust(typ, els)))
	return nil
}

func (w *copyInWalker) copyInHeadersMap(f, t reflect.Value) error {
	if f.Kind() != reflect.String {
		return fmt.Errorf("string expected, got %T", f.Interface())
	}
	h, err := textproto.NewReader(bufio.NewReader(strings.NewReader(f.String()))).ReadMIMEHeader()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	res := make(map[string]attr.Value, len(h))
	for key := range h {
		values := h.Values(key)
		elems := make([]attr.Value, len(values))
		for i, s := range values {
			elems[i] = types.StringValue(s)
		}
		res[key] = types.ListValueMust(types.StringType, elems)
	}
	m := types.MapValueMust(types.ListType{ElemType: types.StringType}, res)
	t.Set(reflect.ValueOf(m))
	return nil
}

func (w *copyInWalker) copyInObject(f, t reflect.Value) error {
	for f.Kind() == reflect.Ptr {
		f = f.Elem()
	}
	if f.Kind() != reflect.Struct {
		return fmt.Errorf("struct expected, got %T", f.Interface())
	}
	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)
	for i := 0; i < f.NumField(); i++ {
		key := stringy.New(f.Type().Field(i).Name).SnakeCase().ToLower()
		ff := f.Field(i)
		switch ff.Kind() {
		case reflect.String:
			attrTypes[key] = types.StringType
			attrValues[key] = types.StringValue(ff.String())
		case reflect.Bool:
			attrTypes[key] = types.BoolType
			attrValues[key] = types.BoolValue(ff.Bool())
		case reflect.Int:
			attrTypes[key] = types.Int64Type
			attrValues[key] = types.Int64Value(ff.Int())
		case reflect.Float64:
			attrTypes[key] = types.Float64Type
			attrValues[key] = types.Float64Value(ff.Float())
		default:
			return fmt.Errorf("%s: unsupported object attribute type: %T", key, ff.Interface())
		}
	}
	t.Set(reflect.ValueOf(types.ObjectValueMust(attrTypes, attrValues)))
	return nil
}

func (w *copyInWalker) copyInNumber(f, t reflect.Value) error {
	for f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return nil
		}
		f = f.Elem()
	}
	dec, ok := f.Interface().(decimal.Decimal)
	if !ok {
		return fmt.Errorf("decimal.Decimal expected, got %T", f.Interface())
	}
	v := types.NumberValue(dec.BigFloat())
	t.Set(reflect.ValueOf(v))
	return nil
}

func (w *copyInWalker) copyInDuration(f, t reflect.Value) error {
	for f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return nil
		}
		f = f.Elem()
	}
	dec, ok := f.Interface().(decimal.Decimal)
	if !ok {
		return fmt.Errorf("decimal.Decimal expected, got %T", f.Interface())
	}
	dec = dec.Mul(decimal.NewFromInt(int64(time.Second)))
	if !dec.IsInteger() {
		return errors.New("resulting duration is not an integer")
	}
	if dec.IsNegative() {
		return errors.New("resulting duration is negative")
	}
	dur := time.Duration(dec.IntPart())
	v := customtypes.NewDurationValue(dur.String())
	t.Set(reflect.ValueOf(v))
	return nil
}

func (w *copyInWalker) copyInSmartPercentage(f, t reflect.Value) error {
	for f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return nil
		}
		f = f.Elem()
	}
	d, ok := f.Interface().(decimal.Decimal)
	if !ok {
		return fmt.Errorf("decimal.Decimal expected, got %T", f.Interface())
	}
	v := customtypes.NewSmartPercentageValue(d.BigFloat())
	t.Set(reflect.ValueOf(v))
	return nil
}
