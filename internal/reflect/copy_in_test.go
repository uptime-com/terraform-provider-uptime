package reflect

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	extratypes "github.com/mikluko/terraform-plugin-framework-extratypes"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCopyIn(t *testing.T) {
	t.Run("types.String", func(t *testing.T) {
		type SrcType struct {
			Foo string
		}
		type DstType struct {
			Foo types.String
		}
		src := SrcType{
			Foo: "Foo",
		}
		dst := DstType{}
		err := CopyIn(&dst, src)
		require.NoError(t, err)
		require.Equal(t, types.StringValue("Foo"), dst.Foo)
	})
	t.Run("types.Bool", func(t *testing.T) {
		type SrcType struct {
			Foo bool
		}
		type DstType struct {
			Foo types.Bool
		}
		src := SrcType{
			Foo: true,
		}
		dst := DstType{}
		err := CopyIn(&dst, src)
		require.NoError(t, err)
		require.Equal(t, types.BoolValue(true), dst.Foo)
	})
	t.Run("types.Int64", func(t *testing.T) {
		type SrcType struct {
			Foo int
		}
		type DstType struct {
			Foo types.Int64
		}
		src := SrcType{
			Foo: 100500,
		}
		dst := DstType{}
		err := CopyIn(&dst, src)
		require.NoError(t, err)
		require.Equal(t, types.Int64Value(100500), dst.Foo)
	})
	t.Run("types.Float64", func(t *testing.T) {
		type SrcType struct {
			Foo float64
		}
		type DstType struct {
			Foo types.Float64
		}
		src := SrcType{
			Foo: 100.500,
		}
		dst := DstType{}
		err := CopyIn(&dst, src)
		require.NoError(t, err)
		require.Equal(t, types.Float64Value(100.500), dst.Foo)
	})
	t.Run("types.List", func(t *testing.T) {
		type SrcType struct {
			Foo []string
		}
		type DstType struct {
			Foo types.List
		}
		src := SrcType{
			Foo: []string{"foo", "bar"},
		}
		dst := DstType{}
		err := CopyIn(&dst, src)
		require.NoError(t, err)
		require.Equal(t, types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
		}), dst.Foo)
	})
	t.Run("types.Set", func(t *testing.T) {
		type SrcType struct {
			Foo []string
		}
		type DstType struct {
			Foo types.Set
		}
		src := SrcType{
			Foo: []string{"foo", "bar"},
		}
		dst := DstType{}
		err := CopyIn(&dst, src)
		require.NoError(t, err)
		require.Equal(t, types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("foo"),
			types.StringValue("bar"),
		}), dst.Foo)
	})
	t.Run("types.Map", func(t *testing.T) {
		t.Run("types.String", func(t *testing.T) {
			type SrcType struct {
				Foo map[string]string
			}
			type DstType struct {
				Foo types.Map
			}
			src := SrcType{
				Foo: map[string]string{"foo": "foo", "bar": "bar"},
			}
			dst := DstType{}
			err := CopyIn(&dst, src)
			require.NoError(t, err)
			require.Equal(t, types.MapValueMust(types.StringType, map[string]attr.Value{
				"foo": types.StringValue("foo"),
				"bar": types.StringValue("bar"),
			}), dst.Foo)
		})
		t.Run("headers", func(t *testing.T) {
			type SrcType struct {
				Foo string
			}
			type DstType struct {
				Foo types.Map `ref:",extra=headers"`
			}
			src := SrcType{
				Foo: "Foo: A\r\nFoo: B\r\nBar: C\r\n",
			}
			dst := DstType{}
			err := CopyIn(&dst, src)
			require.NoError(t, err)
			require.Equal(t,
				types.MapValueMust(
					types.ListType{ElemType: types.StringType},
					map[string]attr.Value{
						"Foo": types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("A"),
							types.StringValue("B"),
						}),
						"Bar": types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("C"),
						}),
					},
				),
				dst.Foo,
			)
		})
	})
	t.Run("types.Object", func(t *testing.T) {
		type SrcTypeFoo struct {
			Bar string
			Baz int
		}
		type SrcType struct {
			Foo SrcTypeFoo
		}
		type DstType struct {
			Foo types.Object
		}
		src := SrcType{
			Foo: SrcTypeFoo{
				Bar: "Bar",
				Baz: 100500,
			},
		}
		dst := DstType{}
		err := CopyIn(&dst, src)
		require.NoError(t, err)
	})
	t.Run("types.Number", func(t *testing.T) {
		dec := decimal.NewFromFloat(100.500)
		type SrcType struct {
			Foo decimal.Decimal
		}
		type DstType struct {
			Foo types.Number
		}
		obj := DstType{}
		err := CopyIn(&obj, SrcType{
			Foo: dec,
		})
		require.NoError(t, err)
		require.True(t, obj.Foo.Equal(types.NumberValue(dec.BigFloat())))
	})
	t.Run("extratypes.Duration", func(t *testing.T) {
		dec := decimal.NewFromFloat(100.500)
		dur := time.Duration(100)*time.Second + time.Duration(500)*time.Millisecond
		type SrcType struct {
			Foo decimal.Decimal
		}
		type DstType struct {
			Foo extratypes.Duration
		}
		obj := DstType{}
		err := CopyIn(&obj, SrcType{
			Foo: dec,
		})
		require.NoError(t, err)
		require.True(t, obj.Foo.Equal(extratypes.NewDurationValue(dur.String())))
	})
	t.Run("options", func(t *testing.T) {
		t.Run("path", func(t *testing.T) {
			type SrcType struct {
				Bar string
			}
			type DstType struct {
				Foo types.String `ref:"Bar"`
			}
			src := SrcType{
				Bar: "Bar",
			}
			dst := DstType{}
			err := CopyIn(&dst, src)
			require.NoError(t, err)
			require.Equal(t, types.StringValue("Bar"), dst.Foo)
		})
		t.Run("opt", func(t *testing.T) {
			type SrcType struct {
			}
			type DstType struct {
				Foo types.String `ref:",opt"`
			}
			src := SrcType{}
			dst := DstType{}
			err := CopyIn(&dst, src)
			require.NoError(t, err)
			require.Equal(t, types.StringNull(), dst.Foo)
		})
		t.Run("skip", func(t *testing.T) {
			type SrcType struct {
				Foo string
			}
			type DstType struct {
				Foo types.String `ref:",skip"`
			}
			src := SrcType{
				Foo: "Foo",
			}
			dst := DstType{}
			err := CopyIn(&dst, src)
			require.NoError(t, err)
			require.Equal(t, types.StringNull(), dst.Foo)
		})
	})
}
