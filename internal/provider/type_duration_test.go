package provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uptime-com/terraform-provider-uptime/internal/provider"
)

func TestDurationImpl(t *testing.T) {
	var (
		_ xattr.TypeWithValidate                     = provider.DurationType
		_ basetypes.StringTypable                    = provider.DurationType
		_ attr.Value                                 = (*provider.Duration)(nil)
		_ basetypes.StringValuableWithSemanticEquals = (*provider.Duration)(nil)
	)
}

func TestDurationType_Validate(t *testing.T) {
	t.Parallel()

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
		"valid duration": {
			in: tftypes.NewValue(tftypes.String, "1h3m10s"),
		},
		"invalid duration": {
			in: tftypes.NewValue(tftypes.String, "1d"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Duration String Value",
					"A string value was provided that is not valid duration.\n\n"+
						"Given Value: 1d\n"+
						"Error: time: unknown unit \"d\" in duration \"1d\"",
				),
			},
		},
		"wrong value type": {
			in: tftypes.NewValue(tftypes.Number, 123),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Duration Type Validation Error",
					"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
						"expected String value, received tftypes.Value with value: tftypes.Number<\"123\">",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := provider.DurationType.Validate(context.Background(), testCase.in, path.Root("test"))

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDurationType_ValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in        tftypes.Value
		expect    attr.Value
		expectErr string
	}{
		"value": {
			in:     tftypes.NewValue(tftypes.String, "1h"),
			expect: provider.DurationValue(time.Hour),
		},
		"unknown": {
			in:     tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expect: provider.DurationUnknown(),
		},
		"null": {
			in:     tftypes.NewValue(tftypes.String, nil),
			expect: provider.DurationNull(),
		},
		"wrong type": {
			in:        tftypes.NewValue(tftypes.Number, 123),
			expectErr: "tftypes.String required, got tftypes.Number",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for name, c := range testCases {
		name, c := name, c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := provider.DurationType.ValueFromTerraform(ctx, c.in)
			if err != nil {
				if c.expectErr == "" {
					t.Fatalf("Unexpected error: %s", err)
				}
				if c.expectErr != err.Error() {
					t.Fatalf("Expected error to be %q, got %q", c.expectErr, err.Error())
				}
				return
			}
			assert.True(t, c.expect.Equal(got))
			assert.True(t, c.expect.IsNull() == c.in.IsNull(), "null-ness mismatch")
			//assert.True(t, c.expect.IsUnknown() == !c.in.IsKnown(), "unknown-ness mismatch")
			if c.expectErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, c.expectErr)
			}
		})
	}
}

func TestDurationValue_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    provider.Duration
		argument provider.Duration
		expected bool
	}{
		"null equals null": {
			value:    provider.DurationNull(),
			argument: provider.DurationNull(),
			expected: true,
		},
		"unknown equals unknown": {
			value:    provider.DurationNull(),
			argument: provider.DurationNull(),
			expected: true,
		},
		"null not equals unknown": {
			value:    provider.DurationNull(),
			argument: provider.DurationUnknown(),
			expected: false,
		},
		"1h equals 1h": {
			value:    provider.DurationValue(time.Hour),
			argument: provider.DurationValue(time.Hour),
			expected: true,
		},
		"1h not equals 1m": {
			value:    provider.DurationValue(time.Hour),
			argument: provider.DurationValue(time.Minute),
			expected: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.Equal(testCase.argument)

			if got != testCase.expected {
				t.Errorf("Unexpected difference in equality, got: %t, expected: %t", got, testCase.expected)
			}
		})
	}
}

func TestDurationValue_StringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    provider.Duration
		argument basetypes.StringValuable
		expected bool
		diags    diag.Diagnostics
	}{
		"null equals null": {
			value:    provider.DurationNull(),
			argument: basetypes.NewStringNull(),
			expected: true,
		},
		"unknown equals unknown": {
			value:    provider.DurationUnknown(),
			argument: basetypes.NewStringUnknown(),
			expected: true,
		},
		"null not equals unknown": {
			value:    provider.DurationNull(),
			argument: basetypes.NewStringUnknown(),
			expected: false,
		},
		"unknown not equals null": {
			value:    provider.DurationUnknown(),
			argument: basetypes.NewStringNull(),
			expected: false,
		},
		"1h equals 1h": {
			value:    provider.DurationValue(time.Hour),
			argument: basetypes.NewStringValue("1h"),
			expected: true,
		},
		"1h equals 60m": {
			value:    provider.DurationValue(time.Hour),
			argument: basetypes.NewStringValue("60m"),
			expected: true,
		},
		"1h not equals 1m": {
			value:    provider.DurationValue(time.Hour),
			argument: basetypes.NewStringValue("1m"),
			expected: false,
		},
		"bad duration string": {
			value:    provider.DurationValue(time.Hour),
			argument: basetypes.NewStringValue("hello, world"),
			expected: false,
			diags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"duration parse error",
					`unexpected error converting string to time.Duration: time: invalid duration "hello, world"`,
				),
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for name, c := range testCases {
		name, c := name, c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, diags := c.value.StringSemanticEquals(ctx, c.argument)
			require.Equal(t, c.expected, got)
			require.Equal(t, c.diags, diags)
		})
	}
}
