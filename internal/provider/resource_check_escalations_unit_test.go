package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestEscalationsToAPI(t *testing.T) {
	testCases := map[string]struct {
		arg    *escalationsAttribute
		expect []upapi.CheckEscalation
	}{
		"nil": {
			arg:    nil,
			expect: nil,
		},
		"empty": {
			arg:    &escalationsAttribute{},
			expect: []upapi.CheckEscalation{},
		},
		"single": {
			arg: &escalationsAttribute{
				{
					WaitTime:      types.Int64Value(300),
					NumRepeats:    types.Int64Value(3),
					ContactGroups: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Default")}),
				},
			},
			expect: []upapi.CheckEscalation{
				{
					WaitTime:      300,
					NumRepeats:    3,
					ContactGroups: &[]string{"Default"},
				},
			},
		},
		"multiple": {
			arg: &escalationsAttribute{
				{
					WaitTime:      types.Int64Value(300),
					NumRepeats:    types.Int64Value(3),
					ContactGroups: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Default")}),
				},
				{
					WaitTime:   types.Int64Value(600),
					NumRepeats: types.Int64Value(0),
					ContactGroups: types.SetValueMust(types.StringType, []attr.Value{
						types.StringValue("Default"),
						types.StringValue("DevOps"),
					}),
				},
			},
			expect: []upapi.CheckEscalation{
				{
					WaitTime:      300,
					NumRepeats:    3,
					ContactGroups: &[]string{"Default"},
				},
				{
					WaitTime:      600,
					NumRepeats:    0,
					ContactGroups: &[]string{"Default", "DevOps"},
				},
			},
		},
	}

	a := *new(escalationsAttributeContextAdapter)
	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			got := a.escalationsToAPI(tc.arg)

			if tc.expect == nil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected %+v, got nil", tc.expect)
			}

			if len(got) != len(tc.expect) {
				t.Fatalf("expected %d escalations, got %d", len(tc.expect), len(got))
			}

			for i := range tc.expect {
				if got[i].WaitTime != tc.expect[i].WaitTime {
					t.Errorf("escalation[%d].WaitTime: expected %d, got %d", i, tc.expect[i].WaitTime, got[i].WaitTime)
				}
				if got[i].NumRepeats != tc.expect[i].NumRepeats {
					t.Errorf("escalation[%d].NumRepeats: expected %d, got %d", i, tc.expect[i].NumRepeats, got[i].NumRepeats)
				}
				gotLen := 0
				expectLen := 0
				if got[i].ContactGroups != nil {
					gotLen = len(*got[i].ContactGroups)
				}
				if tc.expect[i].ContactGroups != nil {
					expectLen = len(*tc.expect[i].ContactGroups)
				}
				if gotLen != expectLen {
					t.Errorf("escalation[%d].ContactGroups: expected %d groups, got %d", i, expectLen, gotLen)
				}
			}
		})
	}
}

func TestEscalationsFromAPI(t *testing.T) {
	testCases := map[string]struct {
		arg    []upapi.CheckEscalation
		expect *escalationsAttribute
	}{
		"nil": {
			arg:    nil,
			expect: nil,
		},
		"empty": {
			arg:    []upapi.CheckEscalation{},
			expect: nil,
		},
		"single": {
			arg: []upapi.CheckEscalation{
				{
					WaitTime:      300,
					NumRepeats:    3,
					ContactGroups: &[]string{"Default"},
				},
			},
			expect: &escalationsAttribute{
				{
					WaitTime:      types.Int64Value(300),
					NumRepeats:    types.Int64Value(3),
					ContactGroups: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Default")}),
				},
			},
		},
		"multiple": {
			arg: []upapi.CheckEscalation{
				{
					WaitTime:      300,
					NumRepeats:    3,
					ContactGroups: &[]string{"Default"},
				},
				{
					WaitTime:      600,
					NumRepeats:    0,
					ContactGroups: &[]string{"Default", "DevOps"},
				},
			},
			expect: &escalationsAttribute{
				{
					WaitTime:      types.Int64Value(300),
					NumRepeats:    types.Int64Value(3),
					ContactGroups: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Default")}),
				},
				{
					WaitTime:   types.Int64Value(600),
					NumRepeats: types.Int64Value(0),
					ContactGroups: types.SetValueMust(types.StringType, []attr.Value{
						types.StringValue("Default"),
						types.StringValue("DevOps"),
					}),
				},
			},
		},
	}

	a := *new(escalationsAttributeContextAdapter)
	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			got := a.escalationsFromAPI(tc.arg)

			if tc.expect == nil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected %d escalations, got nil", len(*tc.expect))
			}

			if len(*got) != len(*tc.expect) {
				t.Fatalf("expected %d escalations, got %d", len(*tc.expect), len(*got))
			}

			for i := range *tc.expect {
				if !(*got)[i].WaitTime.Equal((*tc.expect)[i].WaitTime) {
					t.Errorf("escalation[%d].WaitTime: expected %v, got %v", i, (*tc.expect)[i].WaitTime, (*got)[i].WaitTime)
				}
				if !(*got)[i].NumRepeats.Equal((*tc.expect)[i].NumRepeats) {
					t.Errorf("escalation[%d].NumRepeats: expected %v, got %v", i, (*tc.expect)[i].NumRepeats, (*got)[i].NumRepeats)
				}
				if !(*got)[i].ContactGroups.Equal((*tc.expect)[i].ContactGroups) {
					t.Errorf("escalation[%d].ContactGroups: expected %v, got %v", i, (*tc.expect)[i].ContactGroups, (*got)[i].ContactGroups)
				}
			}
		})
	}
}

func TestEscalationsAttributeContext(t *testing.T) {
	testCases := map[string]struct {
		arg         types.List
		expect      *escalationsAttribute
		expectError bool
	}{
		"null": {
			arg:    types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{}}),
			expect: nil,
		},
		"unknown": {
			arg:    types.ListUnknown(types.ObjectType{AttrTypes: map[string]attr.Type{}}),
			expect: nil,
		},
		"empty": {
			arg: types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"wait_time":      types.Int64Type,
						"num_repeats":    types.Int64Type,
						"contact_groups": types.SetType{ElemType: types.StringType},
					},
				},
				[]attr.Value{},
			),
			expect: &escalationsAttribute{},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := *new(escalationsAttributeContextAdapter)
	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			got, diags := a.escalationsAttributeContext(ctx, tc.arg)

			if tc.expectError && !diags.HasError() {
				t.Error("expected error, got none")
			}

			if tc.expect == nil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected result, got nil")
			}

			if len(*got) != len(*tc.expect) {
				t.Errorf("expected %d escalations, got %d", len(*tc.expect), len(*got))
			}
		})
	}
}
