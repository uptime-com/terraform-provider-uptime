package reflect

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/reflectwalk"
)

func CopyIn(dst any, src any) error {
	return reflectwalk.Walk(dst, &copyInWalker{
		src: src,
	})
}

type copyInWalker struct {
	pathWalker
	src any
	err error
}

func (w *copyInWalker) Exit(loc reflectwalk.Location) error {
	if loc == reflectwalk.Struct && len(w.path) > 0 {
		err := w.handle()
		if err != nil {
			w.err = errors.Join(w.err, err)
		}
	}
	if loc == reflectwalk.WalkLoc && w.err != nil {
		return w.err
	}
	return w.pathWalker.Exit(loc)
}

func (w *copyInWalker) handle() error {
	if w.tag.Skip {
		return nil
	}
	path := w.Path()
	if w.tag.Path != "" {
		path = w.tag.Path
	}
	v, err := FindByPath(w.src, path)
	if err != nil {
		if errors.Is(err, ErrNotFound) && w.tag.Opt {
			return nil
		}
		return err
	}
	switch w.v.Interface().(type) {
	case types.String:
		if v.Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T: %s", v.Interface(), w.Path())
		}
		w.v.Set(reflect.ValueOf(types.StringValue(v.String())))
	case types.Bool:
		if v.Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T: %s", v.Interface(), w.Path())
		}
		w.v.Set(reflect.ValueOf(types.BoolValue(v.Bool())))
	case types.Int64:
		if v.Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T: %s", v.Interface(), w.Path())
		}
		w.v.Set(reflect.ValueOf(types.Int64Value(v.Int())))
	case types.Float64:
		if v.Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T: %s", v.Interface(), w.Path())
		}
		w.v.Set(reflect.ValueOf(types.Float64Value(v.Float())))
	case types.List:
		if v.Kind() != reflect.Slice {
			return fmt.Errorf("slice expected, got %T: %s", v.Interface(), w.Path())
		}
		typ, slc, err := w.typeSlice(v)
		if err != nil {
			return err
		}
		w.v.Set(reflect.ValueOf(types.ListValueMust(typ, slc)))
	case types.Set:
		if v.Kind() != reflect.Slice {
			return fmt.Errorf("slice expected, got %T: %s", v.Interface(), w.Path())
		}
		typ, slc, err := w.typeSlice(v)
		if err != nil {
			return err
		}
		w.v.Set(reflect.ValueOf(types.SetValueMust(typ, slc)))
	case types.Map:
		var (
			typ attr.Type
			m   map[string]attr.Value
		)
		switch w.tag.Extra {
		case "headers":
			if v.Kind() != reflect.String {
				return fmt.Errorf("string expected, got %T: %s", v.Interface(), w.Path())
			}
			typ, m, err = w.typeMapHeaders(v)
		default:
			if v.Kind() != reflect.Map {
				return fmt.Errorf("map expected, got %T: %s", v.Interface(), w.Path())
			}
			typ, m, err = w.typeMap(v)
		}
		if err != nil {
			return err
		}
		w.v.Set(reflect.ValueOf(types.MapValueMust(typ, m)))
	case types.Number:
		return errors.New("not implemented")
	case types.Object:
		return errors.New("not implemented")
	}
	return nil
}

func (w *copyInWalker) typeSlice(v *reflect.Value) (typ attr.Type, slc []attr.Value, err error) {
	switch v.Type().Elem().Kind() {
	case reflect.String:
		typ = types.StringType
		for i := 0; i < v.Len(); i++ {
			slc = append(slc, types.StringValue(v.Index(i).String()))
		}
	case reflect.Bool:
		typ = types.BoolType
		for i := 0; i < v.Len(); i++ {
			slc = append(slc, types.BoolValue(v.Index(i).Bool()))
		}
	case reflect.Int:
		typ = types.Int64Type
		for i := 0; i < v.Len(); i++ {
			slc = append(slc, types.Int64Value(v.Index(i).Int()))
		}
	case reflect.Float64:
		typ = types.Float64Type
		for i := 0; i < v.Len(); i++ {
			slc = append(slc, types.Float64Value(v.Index(i).Float()))
		}
	default:
		err = fmt.Errorf("unsupported slice element type: %T: %s", v.Interface(), w.Path())
	}
	return
}

func (w *copyInWalker) typeMap(v *reflect.Value) (typ attr.Type, m map[string]attr.Value, err error) {
	if v.Type().Key().Kind() != reflect.String {
		return nil, nil, fmt.Errorf("unsupported map key type: %T: %s", v.Interface(), w.Path())
	}
	m = make(map[string]attr.Value)
	switch v.Type().Elem().Kind() {
	case reflect.String:
		typ = types.StringType
		for _, k := range v.MapKeys() {
			m[k.String()] = types.StringValue(v.MapIndex(k).String())
		}
	case reflect.Bool:
		typ = types.BoolType
		for _, k := range v.MapKeys() {
			m[k.String()] = types.BoolValue(v.MapIndex(k).Bool())
		}
	case reflect.Int:
		typ = types.Int64Type
		for _, k := range v.MapKeys() {
			m[k.String()] = types.Int64Value(v.MapIndex(k).Int())
		}
	case reflect.Float64:
		typ = types.Float64Type
		for _, k := range v.MapKeys() {
			m[k.String()] = types.Float64Value(v.MapIndex(k).Float())
		}
	default:
		err = fmt.Errorf("unsupported map element type: %T: %s", v.Interface(), w.Path())
	}
	return typ, m, err
}

func (w *copyInWalker) typeMapHeaders(v *reflect.Value) (attr.Type, map[string]attr.Value, error) {
	header, err := textproto.NewReader(bufio.NewReader(strings.NewReader(v.String()))).ReadMIMEHeader()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, nil, err
	}
	m := make(map[string]attr.Value)
	t := types.ListType{ElemType: types.StringType}
	for key := range header {
		values := header.Values(key)
		elems := make([]attr.Value, 0, len(values))
		for _, s := range values {
			elems = append(elems, types.StringValue(s))
		}
		m[key] = types.ListValueMust(types.StringType, elems)
	}
	return t, m, nil
}
