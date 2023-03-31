package reflect

import (
	"reflect"
	"strings"

	"github.com/mitchellh/reflectwalk"
	"github.com/muir/reflectutils"
)

type pathWalker struct {
	path []string
	tag  tag
	f    reflect.StructField
	v    reflect.Value
}

func (w *pathWalker) Path() string {
	return strings.Join(w.path, ".")
}

func (w *pathWalker) Enter(loc reflectwalk.Location) error {
	if loc == reflectwalk.StructField {
		w.path = append(w.path, w.f.Name)
		w.tag = tag{}
		err := reflectutils.GetTag(w.f.Tag, "ref").Fill(&w.tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *pathWalker) Exit(loc reflectwalk.Location) error {
	if loc == reflectwalk.StructField {
		w.path = w.path[:len(w.path)-1]
	}
	return nil
}

func (w *pathWalker) StructField(f reflect.StructField, _ reflect.Value) error {
	w.f = f
	if !f.IsExported() {
		return reflectwalk.SkipEntry
	}
	return nil
}

func (w *pathWalker) Struct(v reflect.Value) error {
	w.v = v
	return nil
}

type tag struct {
	Path  string `pt:"0"`     // path to field
	Skip  bool   `pt:"skip"`  // skip this field
	Opt   bool   `pt:"opt"`   // this field might be missing on the other side
	Extra string `pt:"extra"` // extra info for the field
}
