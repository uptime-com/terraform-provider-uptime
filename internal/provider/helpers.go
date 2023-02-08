package provider

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ptr[T any](v T) *T {
	return &v
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
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
