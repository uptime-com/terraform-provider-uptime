package provider_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/uptime-com/terraform-provider-uptime/internal/provider"
)

func TestRawJsonType_Validate(t *testing.T) {
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
			in: tftypes.NewValue(tftypes.String, nil),
		},
		"unknown": {
			in: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"valid json": {
			in: tftypes.NewValue(tftypes.String, `["a", "b"]`),
		},
		"invalid json": {
			in: tftypes.NewValue(tftypes.String, `{ ain't a valid json }`),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"RawJsonType Validation Error",
					"Value is not a valid json. This is an error in the configuration.",
				),
			},
		},
		"wrong value type": {
			in: tftypes.NewValue(tftypes.Number, 123),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"RawJsonType Validation Error",
					"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
						"expected String value, received tftypes.Value with value: tftypes.Number<\"123\">",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			diags := new(provider.RawJsonType).Validate(ctx, testCase.in, path.Root("test"))
			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestJsonType_ValueFromTerraform(t *testing.T) {

	testCases := map[string]struct {
		in          tftypes.Value
		expectation attr.Value
		expectedErr string
	}{
		"value": {
			in:          tftypes.NewValue(tftypes.String, `["a", "b"]`),
			expectation: provider.RawJsonValue(`["a", "b"]`),
		},
		"unknown": {
			in:          tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expectation: provider.RawJsonUnknown(),
		},
		"null": {
			in:          tftypes.NewValue(tftypes.String, nil),
			expectation: provider.RawJsonNull(),
		},
		"wrong type": {
			in:          tftypes.NewValue(tftypes.Number, 123),
			expectedErr: "can't unmarshal tftypes.Number into *string, expected string",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			got, err := new(provider.RawJsonType).ValueFromTerraform(ctx, testCase.in)
			if err != nil {
				if testCase.expectedErr == "" {
					t.Fatalf("Unexpected error: %s", err)
				}
				if testCase.expectedErr != err.Error() {
					t.Fatalf("Expected error to be %q, got %q", testCase.expectedErr, err.Error())
				}
				return
			}
			if err == nil && testCase.expectedErr != "" {
				t.Fatalf("Expected error to be %q, didn't get an error", testCase.expectedErr)
			}
			if !got.Equal(testCase.expectation) {
				t.Errorf("Expected %+v, got %+v", testCase.expectation, got)
			}
			if testCase.expectation.IsNull() != testCase.in.IsNull() {
				t.Errorf("Expected null-ness match: expected %t, got %t", testCase.expectation.IsNull(), testCase.in.IsNull())
			}
			if testCase.expectation.IsUnknown() != !testCase.in.IsKnown() {
				t.Errorf("Expected unknown-ness match: expected %t, got %t", testCase.expectation.IsUnknown(), !testCase.in.IsKnown())
			}
		})
	}
}

func TestRawJson_ToTerraformValue(t *testing.T) {
	ctx := context.Background()

	testCases := map[string]struct {
		in       provider.RawJson
		expected tftypes.Value
	}{
		"known value": {
			in:       provider.RawJsonValue(`{"key": "value"}`),
			expected: tftypes.NewValue(tftypes.String, `{"key": "value"}`),
		},
		"unknown value": {
			in:       provider.RawJsonUnknown(),
			expected: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"null value": {
			in:       provider.RawJsonNull(),
			expected: tftypes.NewValue(tftypes.String, nil),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			got, err := testCase.in.ToTerraformValue(ctx)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if !got.Equal(testCase.expected) {
				t.Errorf("Expected %v, got %v", testCase.expected, got)
			}
			if got.IsKnown() != testCase.expected.IsKnown() {
				t.Errorf("Expected IsKnown=%t, got IsKnown=%t", testCase.expected.IsKnown(), got.IsKnown())
			}
			if got.IsNull() != testCase.expected.IsNull() {
				t.Errorf("Expected IsNull=%t, got IsNull=%t", testCase.expected.IsNull(), got.IsNull())
			}
		})
	}
}

func TestRawJson_StringSemanticEquals(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := map[string]struct {
		value1         provider.RawJson
		value2         provider.RawJson
		expectedResult bool
		expectedDiags  diag.Diagnostics
	}{
		"null equals null": {
			value1:         provider.RawJsonNull(),
			value2:         provider.RawJsonNull(),
			expectedResult: true,
		},
		"unknown equals unknown": {
			value1:         provider.RawJsonUnknown(),
			value2:         provider.RawJsonUnknown(),
			expectedResult: true,
		},
		"unknown not equals null": {
			value1:         provider.RawJsonUnknown(),
			value2:         provider.RawJsonNull(),
			expectedResult: false,
		},
		"unknown not equals value": {
			value1:         provider.RawJsonUnknown(),
			value2:         provider.RawJsonValue(`{"a": "b"}`),
			expectedResult: false,
		},
		"null not equals value": {
			value1:         provider.RawJsonNull(),
			value2:         provider.RawJsonValue(`{"a": "b"}`),
			expectedResult: false,
		},
		"value not equals value": {
			value1:         provider.RawJsonValue(`{"a": "b"}`),
			value2:         provider.RawJsonValue(`{"a": "c"}`),
			expectedResult: false,
		},
		"value equals differently formatted value": {
			value1:         provider.RawJsonValue(`{"a": "b"}`),
			value2:         provider.RawJsonValue(`{"a":"b"}\n`),
			expectedResult: true,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			res, diags := testCase.value1.StringSemanticEquals(ctx, testCase.value2)
			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
			if res != testCase.expectedResult {
				t.Errorf("Unexpected result of comparisson, got: %v, expected: %v", res, testCase.expectedResult)
			}
		})
	}
}
