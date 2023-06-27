package reflect

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/reflectwalk"
	"github.com/stretchr/testify/require"
)

type testPathWalker struct {
	path  []string
	tag   []Tag
	value []reflect.Value
}

func (w *testPathWalker) append(p string, t Tag, v reflect.Value) error {
	w.path = append(w.path, p)
	w.tag = append(w.tag, t)
	w.value = append(w.value, v)
	return nil
}

func (w *testPathWalker) Walk(path string, t Tag, v reflect.Value) error {
	return w.append(path, t, v)
}

func TestPathWalker(t *testing.T) {
	t.Run("tree", func(t *testing.T) {
		w := pathWalker{
			PathWalker: &testPathWalker{},
		}
		type Foo struct {
			A string
			B struct {
				C string
				D []string
			}
			E map[string]string
			F types.String
		}
		err := reflectwalk.Walk(Foo{}, &w)
		require.NoError(t, err)
		require.Equal(t, []string{
			"A",
			"B",
			"B.C",
			"B.D",
			"E",
			"F",
		}, w.PathWalker.(*testPathWalker).path)
	})
	t.Run("types", func(t *testing.T) {
		w := pathWalker{
			PathWalker: &testPathWalker{},
		}
		type Foo struct {
			A types.String
			B types.Bool
			C types.Int64
			D types.Float64
		}
		err := reflectwalk.Walk(Foo{}, &w)
		require.NoError(t, err)
		require.Equal(t, []string{
			"A",
			"B",
			"C",
			"D",
		}, w.PathWalker.(*testPathWalker).path)
	})
	t.Run("tag", func(t *testing.T) {
		t.Run("path", func(t *testing.T) {
			w := pathWalker{
				PathWalker: &testPathWalker{},
			}
			type Foo struct {
				A string `ref:"A.B.C"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []Tag{
				{Path: "A.B.C"},
			}, w.PathWalker.(*testPathWalker).tag)
		})
		t.Run("opt", func(t *testing.T) {
			w := pathWalker{
				PathWalker: &testPathWalker{},
			}
			type Foo struct {
				A string `ref:",opt"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []Tag{
				{Opt: true},
			}, w.PathWalker.(*testPathWalker).tag)
		})
		t.Run("skip", func(t *testing.T) {
			w := pathWalker{
				PathWalker: &testPathWalker{},
			}
			type Foo struct {
				A string `ref:",skip"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []Tag{
				{Skip: true},
			}, w.PathWalker.(*testPathWalker).tag)
		})
		t.Run("extra", func(t *testing.T) {
			w := pathWalker{
				PathWalker: &testPathWalker{},
			}
			type Foo struct {
				A string `ref:",extra=bar"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []Tag{
				{Extra: "bar"},
			}, w.PathWalker.(*testPathWalker).tag)
		})
		t.Run("combined", func(t *testing.T) {
			w := pathWalker{
				PathWalker: &testPathWalker{},
			}
			type Foo struct {
				A string `ref:"A.B.C,opt,skip,extra=bar"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []Tag{
				{Path: "A.B.C", Extra: "bar", Opt: true, Skip: true},
			}, w.PathWalker.(*testPathWalker).tag)
		})
	})
}
