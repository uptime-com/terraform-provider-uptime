package reflect

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/reflectwalk"
)

// FindByPath returns the value of the field at the given dot-delimited path.
func FindByPath(obj any, path string) (*reflect.Value, error) {
	w := findWalker{
		target: path,
	}
	err := reflectwalk.Walk(obj, &w)
	return w.res, err
}

var ErrNotFound = errors.New("not found")

type findWalker struct {
	pathWalker
	target string
	cur    *reflect.Value
	res    *reflect.Value
}

func (w *findWalker) Enter(loc reflectwalk.Location) (err error) {
	return w.pathWalker.Enter(loc)
}

func (w *findWalker) Exit(loc reflectwalk.Location) (err error) {
	if w.Path() == w.target {
		w.res = w.cur
	}
	if loc == reflectwalk.WalkLoc {
		if w.res == nil {
			return fmt.Errorf("%w: %s", ErrNotFound, w.target)
		}
	}
	return w.pathWalker.Exit(loc)
}

func (w *findWalker) StructField(f reflect.StructField, v reflect.Value) error {
	w.cur = &v
	return w.pathWalker.StructField(f, v)
}
