package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLegacyOnlyBrandingWarningValidator(t *testing.T) {
	cases := map[string]struct {
		value       types.String
		wantWarning bool
	}{
		"non-empty value warns":   {types.StringValue("<header/>"), true},
		"empty value is silent":   {types.StringValue(""), false},
		"null value is silent":    {types.StringNull(), false},
		"unknown value is silent": {types.StringUnknown(), false},
	}
	v := LegacyOnlyBrandingWarningValidator("custom_header_html_inspire")
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			rq := validator.StringRequest{Path: path.Root("custom_header_html"), ConfigValue: tc.value}
			rs := &validator.StringResponse{}
			v.ValidateString(context.Background(), rq, rs)
			if rs.Diagnostics.HasError() {
				t.Fatalf("validator must never emit an error, got: %v", rs.Diagnostics)
			}
			if got := rs.Diagnostics.WarningsCount() > 0; got != tc.wantWarning {
				t.Fatalf("warning = %v, want %v", got, tc.wantWarning)
			}
			if tc.wantWarning {
				detail := rs.Diagnostics.Warnings()[0].Detail()
				if !strings.Contains(detail, "custom_header_html_inspire") {
					t.Fatalf("warning must point to the *_inspire attribute, got: %q", detail)
				}
			}
		})
	}
}
