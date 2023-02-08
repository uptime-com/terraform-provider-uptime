package provider

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uptime-com/terraform-provider-uptime/internal/uptimeapi"
)

func TestMirror(t *testing.T) {
	t.Run("from api", func(t *testing.T) {
		src := uptimeapi.CheckTag{
			Pk:       ptr(1),
			ColorHex: "#000000",
			Tag:      "example",
			Url:      ptr("https://example.com"),
		}
		type dstType struct {
			Pk       *int
			ColorHex string
			Tag      string
			Url      *string
		}
		dst := dstType{}
		exp := dstType{
			Pk:       ptr(1),
			ColorHex: "#000000",
			Tag:      "example",
			Url:      ptr("https://example.com"),
		}
		mirror(&dst, src)
		require.Equal(t, exp, dst)
	})
	t.Run("to api", func(t *testing.T) {
		type srcType struct {
			Pk       *int
			ColorHex string
			Tag      string
			Url      *string
		}
		src := srcType{
			Pk:       ptr(1),
			ColorHex: "#000000",
			Tag:      "example",
			Url:      ptr("https://example.com"),
		}
		dst := uptimeapi.CheckTag{}
		exp := uptimeapi.CheckTag{
			Pk:       ptr(1),
			ColorHex: "#000000",
			Tag:      "example",
			Url:      ptr("https://example.com"),
		}
		mirror(&dst, src)
		require.Equal(t, exp, dst)
	})
	t.Run("map field dst", func(t *testing.T) {
		src := uptimeapi.CheckTag{
			Url: ptr("https://example.com"),
		}
		type dstType struct {
			URL *string `api:"Url"`
		}
		dst := dstType{}
		exp := dstType{
			URL: ptr("https://example.com"),
		}
		mirror(&dst, src)
		require.Equal(t, exp, dst)
	})
	t.Run("map field src", func(t *testing.T) {
		type srcType struct {
			URL *string `api:"Url"`
		}
		src := srcType{
			URL: ptr("https://example.com"),
		}
		dst := uptimeapi.CheckTag{}
		exp := uptimeapi.CheckTag{
			Url: ptr("https://example.com"),
		}
		mirror(&dst, src)
		require.Equal(t, exp, dst)
	})
	t.Run("map field src", func(t *testing.T) {
		type srcType struct {
			URL *string `api:"Url"`
		}
		src := srcType{
			URL: ptr("https://example.com"),
		}
		dst := uptimeapi.CheckTag{}
		exp := uptimeapi.CheckTag{
			Url: ptr("https://example.com"),
		}
		mirror(&dst, src)
		require.Equal(t, exp, dst)
	})
	t.Run("pk dst", func(t *testing.T) {
		src := uptimeapi.CheckTag{
			Pk: ptr(10),
		}
		type dstType struct {
			ID *string `api:"Pk"`
		}
		dst := dstType{}
		exp := dstType{
			ID: ptr("10"),
		}
		mirror(&dst, src)
		require.Equal(t, exp, dst)
	})
	t.Run("pk src", func(t *testing.T) {
		type srcType struct {
			ID *string `api:"Pk"`
		}
		src := srcType{
			ID: ptr("10"),
		}
		dst := uptimeapi.CheckTag{}
		exp := uptimeapi.CheckTag{
			Pk: ptr(10),
		}
		mirror(&dst, src)
		require.Equal(t, exp, dst)
	})
}
