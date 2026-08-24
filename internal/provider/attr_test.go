package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestStringOptionalAPIValue(t *testing.T) {
	require.Nil(t, stringOptionalAPIValue(types.StringNull()))
	require.Nil(t, stringOptionalAPIValue(types.StringUnknown()))

	empty := stringOptionalAPIValue(types.StringValue(""))
	require.NotNil(t, empty)
	require.Equal(t, "", *empty)

	set := stringOptionalAPIValue(types.StringValue("IPV4"))
	require.NotNil(t, set)
	require.Equal(t, "IPV4", *set)
}

func TestStringOptionalModelValue(t *testing.T) {
	require.True(t, stringOptionalModelValue(nil).IsNull())

	empty := ""
	require.Equal(t, "", stringOptionalModelValue(&empty).ValueString())

	set := "IPV4"
	require.Equal(t, "IPV4", stringOptionalModelValue(&set).ValueString())
}
