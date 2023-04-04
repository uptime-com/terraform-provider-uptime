package reflect

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/reflectwalk"
)

// FindByPath returns the value of the field at the given dot-delimited path.
func FindByPath(obj any, path string) (*reflect.Value, error) {
	f := pathWalkerFind{
		target: path,
	}
	w := pathWalker{
		PathWalker: &f,
	}
	err := reflectwalk.Walk(obj, &w)
	if err != nil {
		return nil, err
	}
	if f.res == nil {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, path)
	}
	return f.res, err
}

var ErrNotFound = errors.New("not found")

type pathWalkerFind struct {
	target string
	res    *reflect.Value
}

func (p *pathWalkerFind) Walk(path string, _ Tag, v reflect.Value) error {
	if path == p.target {
		p.res = &v
	}
	if p.res != nil {
		return reflectwalk.SkipEntry
	}
	return nil
}
