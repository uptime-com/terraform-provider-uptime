package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/uptime-com/terraform-provider-uptime/internal/provider"
)

func TestDecimalImpl(t *testing.T) {
	var (
		_ attr.Type                                  = provider.DecimalType
		_ basetypes.StringTypable                    = provider.DecimalType
		_ attr.Value                                 = (*provider.Decimal)(nil)
		_ basetypes.StringValuableWithSemanticEquals = (*provider.Decimal)(nil)
	)
}
