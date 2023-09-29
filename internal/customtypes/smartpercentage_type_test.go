package customtypes_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"
)

func TestSmartPercentageType_Validate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := map[string]struct {
		in            tftypes.Value
		expectedDiags diag.Diagnostics
	}{
		"empty-struct": {
			in: tftypes.Value{},
		},
		"null": {
			in: tftypes.NewValue(tftypes.Number, nil),
		},
		"unknown": {
			in: tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
		},
		"valid percentage": {
			in: tftypes.NewValue(tftypes.Number, 50.11),
		},
		"invalid percentage precision": {
			in: tftypes.NewValue(tftypes.Number, 50.111),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"SmartPercentage Type Validation Error",
					"Value must have at most 2 digits after decimal point if presented as percentage. This is an error in the configuration.",
				),
			},
		},
		"valid fraction": {
			in: tftypes.NewValue(tftypes.Number, 0.9999),
		},
		"invalid fraction precision": {
			in: tftypes.NewValue(tftypes.Number, 0.99999),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"SmartPercentage Type Validation Error",
					"Value must have at most 4 digits after decimal point if presented as fraction. This is an error in the configuration.",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := new(customtypes.SmartPercentageType).Validate(ctx, testCase.in, path.Root("test"))

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestSmartPercentageType_ValueFromTerraform(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := map[string]struct {
		in        tftypes.Value
		expect    attr.Value
		expectErr string
	}{
		"unknown": {
			in:     tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			expect: customtypes.NewSmartPercentageUnknown(),
		},
		"null": {
			in:     tftypes.NewValue(tftypes.Number, nil),
			expect: customtypes.NewSmartPercentageNull(),
		},
		"wrong type": {
			in:        tftypes.NewValue(tftypes.String, "hello"),
			expectErr: "can't unmarshal tftypes.String into *big.Float, expected *big.Float",
		},
		"percentage": {
			in:     tftypes.NewValue(tftypes.Number, big.NewFloat(50.111)), // it should be ok to have more precision than 2 digits at this point
			expect: customtypes.NewSmartPercentageValue(big.NewFloat(50.111)),
		},
		"fraction": {
			in:     tftypes.NewValue(tftypes.Number, big.NewFloat(0.50111)), // it should be ok to have more precision than 4 digits at this point
			expect: customtypes.NewSmartPercentageValueFraction(big.NewFloat(0.50111)),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := new(customtypes.SmartPercentageType).ValueFromTerraform(ctx, testCase.in)
			if err != nil {
				if testCase.expectErr == "" {
					t.Fatalf("Unexpected error: %s", err)
				}
				if testCase.expectErr != err.Error() {
					t.Fatalf("Expected error to be %q, got %q", testCase.expectErr, err.Error())
				}
				return
			}
			if err == nil && testCase.expectErr != "" {
				t.Fatalf("Expected error to be %q, didn't get an error", testCase.expectErr)
			}
			if !got.Equal(testCase.expect) {
				t.Errorf("Expected %+v, got %+v", testCase.expect, got)
			}
			if testCase.expect.IsNull() != testCase.in.IsNull() {
				t.Errorf("Expected null-ness match: expected %t, got %t", testCase.expect.IsNull(), testCase.in.IsNull())
			}
			if testCase.expect.IsUnknown() != !testCase.in.IsKnown() {
				t.Errorf("Expected unknown-ness match: expected %t, got %t", testCase.expect.IsUnknown(), !testCase.in.IsKnown())
			}
		})
	}
}
