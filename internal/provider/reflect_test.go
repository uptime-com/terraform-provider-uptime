package provider

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromTerraform(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("values", func(t *testing.T) {
		type srcType struct {
			Bool types.Bool   `map:"Bool"`
			Str  types.String `map:"Str"`
			Int  types.Int64  `map:"Int"`
			Num  types.Number `map:"Num"`
			Set  types.Set    `map:"Set"`
		}
		type dstType struct {
			Bool bool
			Str  string
			Int  int
			Num  float64
			Set  []string
		}
		src := srcType{
			Bool: basetypes.NewBoolValue(true),
			Str:  basetypes.NewStringValue("foo"),
			Int:  basetypes.NewInt64Value(10),
			Num:  basetypes.NewNumberValue(big.NewFloat(10.5)),
			Set: mustDiag(basetypes.NewSetValue(basetypes.StringType{}, []attr.Value{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			})),
		}
		dst := dstType{}
		exp := dstType{
			Bool: true,
			Str:  "foo",
			Int:  10,
			Num:  10.5,
			Set:  []string{"a", "b", "c"},
		}
		d := fromTerraform(ctx, &dst, src)
		require.False(t, d.HasError())
		assert.Equal(t, exp, dst)
	})
	t.Run("pointers", func(t *testing.T) {
		type srcType struct {
			Bool types.Bool   `map:"Bool"`
			Str  types.String `map:"Str"`
			Int  types.Int64  `map:"Int"`
			Num  types.Number `map:"Num"`
			Set  types.Set    `map:"Set"`
		}
		type dstType struct {
			Bool *bool
			Str  *string
			Int  *int
			Num  *float64
			Set  *[]string
		}
		src := srcType{
			Bool: basetypes.NewBoolValue(true),
			Str:  basetypes.NewStringValue("foo"),
			Int:  basetypes.NewInt64Value(10),
			Num:  basetypes.NewNumberValue(big.NewFloat(10.5)),
			Set: mustDiag(basetypes.NewSetValue(basetypes.StringType{}, []attr.Value{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			})),
		}
		dst := dstType{}
		exp := dstType{
			Bool: ptr(true),
			Str:  ptr("foo"),
			Int:  ptr(10),
			Num:  ptr(10.5),
			Set:  ptr([]string{"a", "b", "c"}),
		}
		d := fromTerraform(ctx, &dst, src)
		require.False(t, d.HasError())
		assert.Equal(t, exp, dst)
	})

	t.Run("null value", func(t *testing.T) {
		type srcType struct {
			Bool  types.Bool `map:"Bool"`
			BoolP types.Bool `map:"BoolP"`
		}
		type dstType struct {
			Bool  bool
			BoolP *bool
		}
		src := srcType{}
		dst := dstType{}
		exp := dstType{
			Bool:  false,
			BoolP: nil,
		}
		d := fromTerraform(ctx, &dst, src)
		require.True(t, d.HasError())
		assert.Len(t, d, 1)
		assert.Equal(t, "cannot set null value to non-pointer field", d[0].Summary())
		assert.Equal(t, exp, dst)
	})
}

func TestToTerraform(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		type dstType struct {
			Bool types.Bool   `map:"Bool"`
			Str  types.String `map:"Str"`
			Int  types.Int64  `map:"Int"`
			Set  types.Set    `map:"Set"`
		}
		type srcType struct {
			Bool bool
			Str  string
			Int  int64
			Set  []string
		}
		src := srcType{
			Bool: true,
			Str:  "foo",
			Int:  10,
			Set:  []string{"a", "b", "c"},
		}
		dst := dstType{}
		d := toTerraform(&dst, src)
		if d.HasError() {
			t.Error(d)
		}
		exp := dstType{
			Bool: basetypes.NewBoolValue(true),
			Str:  basetypes.NewStringValue("foo"),
			Int:  basetypes.NewInt64Value(10),
			Set: mustDiag(basetypes.NewSetValue(basetypes.StringType{}, []attr.Value{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			})),
		}
		require.Equal(t, exp, dst)
	})
	t.Run("pointers", func(t *testing.T) {
		type dstStruct struct {
			Bool types.Bool   `map:"Bool"`
			Str  types.String `map:"Str"`
			Int  types.Int64  `map:"Int"`
			Set  types.Set    `map:"Set"`
		}
		type srcStruct struct {
			Bool *bool
			Str  *string
			Int  *int64
			Set  *[]string
		}
		src := srcStruct{
			Bool: ptr(true),
			Str:  ptr("foo"),
			Int:  ptr(int64(10)),
			Set:  &[]string{"a", "b", "c"},
		}
		dst := dstStruct{}
		d := toTerraform(&dst, src)
		if d.HasError() {
			t.Error(d)
		}
		exp := dstStruct{
			Bool: basetypes.NewBoolValue(true),
			Str:  basetypes.NewStringValue("foo"),
			Int:  basetypes.NewInt64Value(10),
			Set: mustDiag(basetypes.NewSetValue(basetypes.StringType{}, []attr.Value{
				basetypes.NewStringValue("a"),
				basetypes.NewStringValue("b"),
				basetypes.NewStringValue("c"),
			})),
		}
		require.Equal(t, exp, dst)
	})
	t.Run("unsupported type", func(t *testing.T) {
		type dstStruct struct {
			Unsupported types.Bool `map:"Unsupported"`
		}
		type srcStruct struct {
			Unsupported time.Time
		}
		src := srcStruct{
			Unsupported: time.Now(),
		}
		dst := dstStruct{}
		d := toTerraform(&dst, src)
		require.True(t, d.HasError())
		require.Contains(t, d[0].Summary(), "unsupported type")
	})
}
