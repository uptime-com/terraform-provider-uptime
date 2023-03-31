package reflect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	type Foo struct {
		B string
		C struct {
			D string
		}
		E []string
	}
	foo := Foo{
		B: "hello",
		C: struct{ D string }{D: "world"},
		E: []string{"a", "b", "c"},
	}
	t.Run("top level", func(t *testing.T) {
		res, err := FindByPath(foo, "B")
		require.NoError(t, err)
		assert.Equal(t, "hello", res.String())
	})
	t.Run("nested", func(t *testing.T) {
		res, err := FindByPath(foo, "C.D")
		require.NoError(t, err)
		assert.Equal(t, "world", res.String())
	})
	t.Run("not found", func(t *testing.T) {
		res, err := FindByPath(foo, "Z.Z")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, res)
	})
	t.Run("kind-slice", func(t *testing.T) {
		res, err := FindByPath([]Foo{foo}, "E")
		require.NoError(t, err)
		x := res.Interface()
		assert.Equal(t, []string{"a", "b", "c"}, x.([]string))
	})
}
