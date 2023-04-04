package reflect

import (
	"errors"
	"reflect"
	"strings"

	"github.com/mitchellh/reflectwalk"
	"github.com/muir/reflectutils"
)

type PathWalker interface {
	Walk(path string, t Tag, v reflect.Value) error
}

type pathWalker struct {
	PathWalker
	path []reflect.StructField
}

func (w *pathWalker) Path() string {
	path := make([]string, len(w.path))
	for i, f := range w.path {
		path[i] = f.Name
	}
	return strings.Join(path, ".")
}

func (w *pathWalker) Enter(loc reflectwalk.Location) error {
	return nil
}

func (w *pathWalker) Exit(loc reflectwalk.Location) error {
	if loc == reflectwalk.StructField {
		w.path = w.path[:len(w.path)-1]
	}
	return nil
}

func (w *pathWalker) Struct(v reflect.Value) error {
	return nil
}

func (w *pathWalker) StructField(f reflect.StructField, v reflect.Value) error {
	if f.Anonymous {
		return reflectwalk.SkipEntry
	}
	if !f.IsExported() {
		return reflectwalk.SkipEntry
	}

	w.path = append(w.path, f)

	var tag Tag
	if err := reflectutils.GetTag(f.Tag, "ref").Fill(&tag); err != nil {
		return err
	}
	return w.Walk(w.Path(), tag, v)
}

func (w *pathWalker) Walk(p string, t Tag, v reflect.Value) error {
	if w.PathWalker == nil {
		return errors.New("no path walker")
	}
	return w.PathWalker.Walk(p, t, v)
}

type Tag struct {
	Path  string `pt:"0"`     // path to field
	Skip  bool   `pt:"skip"`  // skip this field
	Opt   bool   `pt:"opt"`   // this field might be missing on the other side
	Extra string `pt:"extra"` // extra info for the field
}
