package customtypes_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uptime-com/terraform-provider-uptime/internal/customtypes"
)

func TestSmartPercentage_Constructors(t *testing.T) {
	testCases := map[string]struct {
		value  customtypes.SmartPercentage
		expect *big.Float
	}{
		"null": {
			value:  customtypes.NewSmartPercentageNull(),
			expect: big.NewFloat(0),
		},
		"unknown": {
			value:  customtypes.NewSmartPercentageUnknown(),
			expect: big.NewFloat(0),
		},
		"percentage": {
			value:  customtypes.NewSmartPercentageValue(big.NewFloat(12.345)),
			expect: big.NewFloat(0.12345),
		},
		"fraction": {
			value:  customtypes.NewSmartPercentageValueFraction(big.NewFloat(0.12345)),
			expect: big.NewFloat(0.12345),
		},
		"smart percentage": {
			value:  customtypes.NewSmartPercentageValue(big.NewFloat(12.345)),
			expect: big.NewFloat(0.12345),
		},
		"smart fraction": {
			value:  customtypes.NewSmartPercentageValue(big.NewFloat(0.12345)),
			expect: big.NewFloat(0.12345),
		},
		"1 equals 1 percent": {
			value:  customtypes.NewSmartPercentageValue(big.NewFloat(1.0)),
			expect: big.NewFloat(0.01),
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			require.Equal(t, 0, testCase.expect.Cmp(testCase.value.ValueBigFloat()), "expected %s, got %s", testCase.expect.String(), testCase.value.ValueBigFloat().String())
		})
	}
}

func TestSmartPercentage_Equal(t *testing.T) {
	testCases := map[string]struct {
		value    customtypes.SmartPercentage
		argument customtypes.SmartPercentage
		expected bool
	}{
		"null equals null": {
			value:    customtypes.NewSmartPercentageNull(),
			argument: customtypes.NewSmartPercentageNull(),
			expected: true,
		},
		"unknown equals unknown": {
			value:    customtypes.NewSmartPercentageNull(),
			argument: customtypes.NewSmartPercentageNull(),
			expected: true,
		},
		"null not equals unknown": {
			value:    customtypes.NewSmartPercentageNull(),
			argument: customtypes.NewSmartPercentageUnknown(),
			expected: false,
		},
		"0.5 equals 50": {
			value:    customtypes.NewSmartPercentageValue(big.NewFloat(50)),
			argument: customtypes.NewSmartPercentageValueFraction(big.NewFloat(0.5)),
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			require.Equal(t, testCase.expected, testCase.value.Equal(testCase.argument))
		})
	}
}
