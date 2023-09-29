package customtypes_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"
)

func TestDuration_ValueDuration(t *testing.T) {
	t.Parallel()

	mustParseDuration := func(s string) time.Duration {
		d, err := time.ParseDuration(s)
		if err != nil {
			t.Fatalf("failed to parse duration %q: %s", s, err)
		}
		return d
	}

	testCases := map[string]struct {
		value         customtypes.Duration
		expectedValue time.Duration
		expectedDiags diag.Diagnostics
	}{
		"value is null": {
			value: customtypes.NewDurationNull(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Duration ValueDuration Error",
					"duration string value is null",
				),
			},
		},
		"value is unknown": {
			value: customtypes.NewDurationUnknown(),
			expectedDiags: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Duration ValueDuration Error",
					"duration string value is unknown",
				),
			},
		},
		"valid integer": {
			value:         customtypes.NewDurationValue("1h"),
			expectedValue: mustParseDuration("1h"),
		},
		"valid float": {
			value:         customtypes.NewDurationValue("1.5s"),
			expectedValue: mustParseDuration("1500ms"),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dur, diags := testCase.value.ValueDuration()

			if dur != testCase.expectedValue {
				t.Errorf("Unexpected difference in time.Duration, got: %s, expected: %s", dur, testCase.expectedValue)
			}

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDuration_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    customtypes.Duration
		argument customtypes.Duration
		expected bool
	}{
		"null equals null": {
			value:    customtypes.NewDurationNull(),
			argument: customtypes.NewDurationNull(),
			expected: true,
		},
		"unknown equals unknown": {
			value:    customtypes.NewDurationNull(),
			argument: customtypes.NewDurationNull(),
			expected: true,
		},
		"null not equals unknown": {
			value:    customtypes.NewDurationNull(),
			argument: customtypes.NewDurationUnknown(),
			expected: false,
		},
		"1h equals 1h": {
			value:    customtypes.NewDurationValue("1h"),
			argument: customtypes.NewDurationValue("1h"),
			expected: true,
		},
		"1h30m equals 560s": {
			value:    customtypes.NewDurationValue("1h"),
			argument: customtypes.NewDurationValue("60m"),
			expected: true,
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
