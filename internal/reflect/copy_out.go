package reflect

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/reflectwalk"
)

func CopyOut(dst any, src any) error {
	return reflectwalk.Walk(src, &copyOutWalker{
		dst: dst,
	})
}

type copyOutWalker struct {
	pathWalker
	dst any
	err error
}

func (w *copyOutWalker) Exit(loc reflectwalk.Location) error {
	if loc == reflectwalk.Struct && w.v.CanInterface() && len(w.path) > 0 {
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

func (w *copyOutWalker) handle() error {
	if a, ok := w.v.Interface().(attr.Value); !ok {
		return nil // skip value which isn't an attribute
	} else if a.IsUnknown() {
		return nil // skip value without value set
	}
	if w.tag.Skip {
		return nil
	}
	path := w.Path()
	if w.tag.Path != "" {
		path = w.tag.Path
	}
	v, err := FindByPath(w.dst, path)
	if err != nil {
		if errors.Is(err, ErrNotFound) && w.tag.Opt {
			return nil
		}
		return err
	}
	if !v.CanSet() {
		return fmt.Errorf("value is not settable: %s", w.Path())
	}
	switch x := w.v.Interface().(type) {
	case types.String:
		if v.Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T: %s", v.Interface(), w.Path())
		}
		v.SetString(x.ValueString())
	case types.Bool:
		if v.Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T: %s", v.Interface(), w.Path())
		}
		v.SetBool(x.ValueBool())
	case types.Int64:
		if v.Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T: %s", v.Interface(), w.Path())
		}
		v.SetInt(x.ValueInt64())
	case types.Float64:
		if v.Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T: %s", v.Interface(), w.Path())
		}
		v.SetFloat(x.ValueFloat64())
	case types.List:
		if v.Kind() != reflect.Slice {
			return fmt.Errorf("slice expected, got %T: %s", v.Interface(), w.Path())
		}
		return w.sliceTo(x, v)
	case types.Set:
		if v.Kind() != reflect.Slice {
			return fmt.Errorf("slice expected, got %T: %s", v.Interface(), w.Path())
		}
		return w.sliceTo(x, v)
	case types.Map:
		if v.Kind() != reflect.Map {
			return fmt.Errorf("map expected, got %T: %s", v.Interface(), w.Path())
		}
		return w.mapTo(x, v)
	case types.Number:
		return errors.New("not implemented")
	case types.Object:
		return errors.New("not implemented")
	}
	return nil
}

type sliceElementable interface {
	Elements() []attr.Value
	ElementType(ctx context.Context) attr.Type
}

func (w *copyOutWalker) sliceTo(x sliceElementable, v *reflect.Value) error {
	els := x.Elements()
	v.Grow(len(els))
	v.SetLen(len(els))
	switch et := x.ElementType(nil); et {
	case types.StringType:
		if v.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T: %s", v.Interface(), w.Path())
		}
		for i, e := range x.Elements() {
			v.Index(i).Set(reflect.ValueOf(e.(types.String).ValueString()))
		}
	case types.BoolType:
		if v.Type().Elem().Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T: %s", v.Interface(), w.Path())
		}
	case types.Int64Type:
		if v.Type().Elem().Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T: %s", v.Interface(), w.Path())
		}
	case types.Float64Type:
		if v.Type().Elem().Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T: %s", v.Interface(), w.Path())
		}
	default:
		return fmt.Errorf("unsupported slice or set element type: %T, %s", et, w.Path())
	}
	return nil
}

type mapElementable interface {
	Elements() map[string]attr.Value
	ElementType(ctx context.Context) attr.Type
}

func (w *copyOutWalker) mapTo(x mapElementable, v *reflect.Value) error {
	els := x.Elements()
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	switch et := x.ElementType(nil); et {
	case types.StringType:
		if v.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("string expected, got %T: %s", v.Interface(), w.Path())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.String).ValueString()))
		}
	case types.BoolType:
		if v.Type().Elem().Kind() != reflect.Bool {
			return fmt.Errorf("bool expected, got %T: %s", v.Interface(), w.Path())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.Bool).ValueBool()))
		}
	case types.Int64Type:
		if v.Type().Elem().Kind() != reflect.Int {
			return fmt.Errorf("int expected, got %T: %s", v.Interface(), w.Path())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.Int64).ValueInt64()))
		}
	case types.Float64Type:
		if v.Type().Elem().Kind() != reflect.Float64 {
			return fmt.Errorf("float64 expected, got %T: %s", v.Interface(), w.Path())
		}
		for key, e := range els {
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(e.(types.Float64).ValueFloat64()))
		}
	default:
		return fmt.Errorf("unsupported map element type: %T, %s", et, w.Path())
	}
	return nil
}
