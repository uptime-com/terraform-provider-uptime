package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestHeadersAttributeAdapter_HeadersAttributeContext(t *testing.T) {
	testCases := map[string]struct {
		arg         types.Map
		expect      string
		expectError bool
	}{
		"simple": {
			arg: types.MapValueMust(HeadersType.ElementType(), map[string]attr.Value{
				"Destination": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Eschaton")}),
			}),
			expect: "Destination: Eschaton\r\n",
		},
		"multiple": {
			arg: types.MapValueMust(HeadersType.ElementType(), map[string]attr.Value{
				"Foo": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Bar"), types.StringValue("Baz")}),
				"Qux": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Quux")}),
			}),
			expect: "Foo: Bar\r\n" +
				"Foo: Baz\r\n" +
				"Qux: Quux\r\n",
		},
		"null": {
			arg: types.MapNull(HeadersType.ElementType()),
		},
		"unknown": {
			arg: types.MapUnknown(HeadersType.ElementType()),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := *new(HeadersAttributeAdapter)
	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			got, diags := a.HeadersAttributeContext(ctx, tc.arg)
			if got != tc.expect {
				t.Errorf("expected %q, got %q", tc.expect, got)
			}
			if tc.expectError && !diags.HasError() {
				t.Error("expected error, got none")
			}
		})
	}
}

func TestHeadersAttributeAdapter_HeadersAttributeValue(t *testing.T) {
	testCases := map[string]struct {
		arg         string
		expect      types.Map
		expectError bool
	}{
		"simple": {
			arg: "Destination: Eschaton\r\n",
			expect: types.MapValueMust(HeadersType.ElementType(), map[string]attr.Value{
				"Destination": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Eschaton")}),
			}),
		},
		"multiple": {
			arg: "Foo: Bar\r\n" +
				"Foo: Baz\r\n" +
				"Qux: Quux\r\n",
			expect: types.MapValueMust(HeadersType.ElementType(), map[string]attr.Value{
				"Foo": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Bar"), types.StringValue("Baz")}),
				"Qux": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Quux")}),
			}),
		},
		"empty": {
			arg:    "",
			expect: types.MapValueMust(HeadersType.ElementType(), map[string]attr.Value{}),
		},
	}

	a := *new(HeadersAttributeAdapter)
	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			got, err := a.HeadersAttributeValue(tc.arg)
			if tc.expectError && err == nil {
				t.Fatal("expected error, got none")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.expect.Equal(got) {
				t.Errorf("expected %s, got %s", tc.expect, got)
				return
			}
		})
	}
}
