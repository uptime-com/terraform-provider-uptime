package customtypes_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"
)

func TestJson_StringSemanticEquals(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := map[string]struct {
		value1         customtypes.Json
		value2         customtypes.Json
		expectedResult bool
		expectedDiags  diag.Diagnostics
	}{
		"null equals null": {
			value1:         customtypes.NewJsonNull(),
			value2:         customtypes.NewJsonNull(),
			expectedResult: true,
		},
		"unknown equals unknown": {
			value1:         customtypes.NewJsonUnknown(),
			value2:         customtypes.NewJsonUnknown(),
			expectedResult: true,
		},
		"unknown equals null": {
			value1:         customtypes.NewJsonUnknown(),
			value2:         customtypes.NewJsonNull(),
			expectedResult: true,
		},
		"unknown equals value": {
			value1:         customtypes.NewJsonUnknown(),
			value2:         customtypes.NewJsonValue(`{"a": "b"}`),
			expectedResult: true,
		},
		"null not equals value": {
			value1:         customtypes.NewJsonNull(),
			value2:         customtypes.NewJsonValue(`{"a": "b"}`),
			expectedResult: false,
		},
		"value not equals value": {
			value1:         customtypes.NewJsonValue(`{"a": "b"}`),
			value2:         customtypes.NewJsonValue(`{"a": "c"}`),
			expectedResult: false,
		},
		"value equals differently formatted value": {
			value1:         customtypes.NewJsonValue(`{"a": "b"}`),
			value2:         customtypes.NewJsonValue(`{"a":"b"}\n`),
			expectedResult: true,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
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
