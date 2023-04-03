package reflect

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCopyOut(t *testing.T) {
	t.Run("types.String", func(t *testing.T) {
		type SrcType struct {
			Foo types.String
		}
		type DstType struct {
			Foo string
		}
		src := SrcType{
			Foo: types.StringValue("Foo"),
		}
		dst := DstType{}
		err := CopyOut(&dst, src)
		require.NoError(t, err)
		require.Equal(t, "Foo", dst.Foo)
	})
	t.Run("types.Bool", func(t *testing.T) {
		type SrcType struct {
			Foo types.Bool
		}
		type DstType struct {
			Foo bool
		}
		src := SrcType{
			Foo: types.BoolValue(true),
		}
		dst := DstType{}
		err := CopyOut(&dst, src)
		require.NoError(t, err)
		require.Equal(t, true, dst.Foo)
	})
	t.Run("types.Int64", func(t *testing.T) {
		type SrcType struct {
			Foo types.Int64
		}
		type DstType struct {
			Foo int
		}
		src := SrcType{
			Foo: types.Int64Value(100500),
		}
		dst := DstType{}
		err := CopyOut(&dst, src)
		require.NoError(t, err)
		require.Equal(t, 100500, dst.Foo)
	})
	t.Run("types.Float64", func(t *testing.T) {
		type SrcType struct {
			Foo types.Float64
		}
		type DstType struct {
			Foo float64
		}
		src := SrcType{
			Foo: types.Float64Value(100.500),
		}
		dst := DstType{}
		err := CopyOut(&dst, src)
		require.NoError(t, err)
		require.Equal(t, 100.500, dst.Foo)
	})
	t.Run("types.List", func(t *testing.T) {
		type SrcType struct {
			Foo types.List
		}
		type DstType struct {
			Foo []string
		}
		src := SrcType{
			Foo: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("Foo"),
				types.StringValue("Bar"),
			}),
		}
		dst := DstType{}
		err := CopyOut(&dst, &src)
		require.NoError(t, err)
		require.Equal(t, []string{"Foo", "Bar"}, dst.Foo)
	})
	t.Run("types.Set", func(t *testing.T) {
		type SrcType struct {
			Foo types.Set
		}
		type DstType struct {
			Foo []string
		}
		src := SrcType{
			Foo: types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("Foo"),
				types.StringValue("Bar"),
			}),
		}
		dst := DstType{}
		err := CopyOut(&dst, &src)
		require.NoError(t, err)
		require.Equal(t, []string{"Foo", "Bar"}, dst.Foo)
	})
	t.Run("types.Map", func(t *testing.T) {
		type SrcType struct {
			Foo types.Map
		}
		type DstType struct {
			Foo map[string]string
		}
		src := SrcType{
			Foo: types.MapValueMust(types.StringType, map[string]attr.Value{
				"Foo": types.StringValue("Foo"),
				"Bar": types.StringValue("Bar"),
			}),
		}
		dst := DstType{}
		err := CopyOut(&dst, &src)
		require.NoError(t, err)
		require.Equal(t, map[string]string{"Foo": "Foo", "Bar": "Bar"}, dst.Foo)
	})
	t.Run("options", func(t *testing.T) {
		t.Run("path", func(t *testing.T) {
			type SrcType struct {
				Foo types.String `ref:"Bar"`
			}
			type DstType struct {
				Bar string
			}
			src := SrcType{
				Foo: types.StringValue("Foo"),
			}
			dst := DstType{}
			err := CopyOut(&dst, src)
			require.NoError(t, err)
			require.Equal(t, "Foo", dst.Bar)
		})
		t.Run("opt", func(t *testing.T) {
			type SrcType struct {
				Foo types.String `ref:",opt"`
			}
			type DstType struct {
			}
			src := SrcType{}
			dst := DstType{}
			err := CopyOut(&dst, src)
			require.NoError(t, err)
		})
		t.Run("skip", func(t *testing.T) {
			type SrcType struct {
				Foo types.String `ref:",skip"`
			}
			type DstType struct {
				Foo string
			}
			src := SrcType{
				Foo: types.StringValue("Foo"),
			}
			dst := DstType{}
			err := CopyOut(&dst, src)
			require.NoError(t, err)
			require.Equal(t, "", dst.Foo)
		})
		t.Run("extra", func(t *testing.T) {
			t.Run("headers", func(t *testing.T) {
				type SrcType struct {
					Foo types.Map `ref:",extra=headers"`
				}
				type DstType struct {
					Foo string
				}
				src := SrcType{
					Foo: types.MapValueMust(
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
				}
				dst := DstType{}
				err := CopyOut(&dst, src)
				require.NoError(t, err)
				require.Contains(t, dst.Foo, "Foo: A\r\n")
				require.Contains(t, dst.Foo, "Foo: B\r\n")
				require.Contains(t, dst.Foo, "Bar: C\r\n")
			})
		})

	})
}
