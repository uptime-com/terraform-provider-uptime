package uptime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaderConversion(t *testing.T) {
	hm := make(map[string]interface{})
	hm["Foo"] = "bar"
	hm["Baz"] = "bat"

	require.Contains(t, headersMapToString(hm), "Foo: bar")
	require.Contains(t, headersMapToString(hm), "Baz: bat")
}
