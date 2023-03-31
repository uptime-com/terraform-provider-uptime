package reflect

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/reflectwalk"
	"github.com/stretchr/testify/require"
)

type testPathWalker struct {
	pathWalker
	pathHistory [][]string
	tagHistory  []tag
}

func (w *testPathWalker) Exit(loc reflectwalk.Location) error {
	if loc == reflectwalk.StructField {
		w.pathHistory = append(w.pathHistory, make([]string, len(w.path)))
		copy(w.pathHistory[len(w.pathHistory)-1], w.path)
		w.tagHistory = append(w.tagHistory, w.tag)
	}
	return w.pathWalker.Exit(loc)
}

func TestPathWalker(t *testing.T) {
	t.Run("tree", func(t *testing.T) {
		w := testPathWalker{}
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
		require.Equal(t, [][]string{
			{"A"},
			{"B", "C"},
			{"B", "D"},
			{"B"},
			{"E"},
			{"F"},
		}, w.pathHistory)
	})
	t.Run("tag", func(t *testing.T) {
		t.Run("path", func(t *testing.T) {
			w := testPathWalker{}
			type Foo struct {
				A string `ref:"A.B.C"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []tag{
				{Path: "A.B.C"},
			}, w.tagHistory)
		})
		t.Run("opt", func(t *testing.T) {
			w := testPathWalker{}
			type Foo struct {
				A string `ref:",opt"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []tag{
				{Opt: true},
			}, w.tagHistory)
		})
		t.Run("skip", func(t *testing.T) {
			w := testPathWalker{}
			type Foo struct {
				A string `ref:",skip"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []tag{
				{Skip: true},
			}, w.tagHistory)
		})
		t.Run("extra", func(t *testing.T) {
			w := testPathWalker{}
			type Foo struct {
				A string `ref:",extra=bar"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []tag{
				{Extra: "bar"},
			}, w.tagHistory)
		})
		t.Run("combined", func(t *testing.T) {
			w := testPathWalker{}
			type Foo struct {
				A string `ref:"A.B.C,opt,skip,extra=bar"`
			}
			err := reflectwalk.Walk(Foo{}, &w)
			require.NoError(t, err)
			require.Equal(t, []tag{
				{Path: "A.B.C", Extra: "bar", Opt: true, Skip: true},
			}, w.tagHistory)
		})
	})
}
