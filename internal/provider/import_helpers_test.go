package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TestImportStateCompositeIDForRejectsBadID verifies the handler surfaces a malformed
// import ID as a diagnostic naming the resource's own parent attribute, rather than the
// statuspage_id the shared parser used to hardcode.
func TestImportStateCompositeIDForRejectsBadID(t *testing.T) {
	handler := ImportStateCompositeIDFor("service_id")
	resp := &resource.ImportStateResponse{}
	handler(context.Background(), resource.ImportStateRequest{ID: "41216"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for an import ID with no parent half")
	}
	if detail := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(detail, "service_id:resource_id") {
		t.Errorf("diagnostic should name the parent attribute, got: %s", detail)
	}
}

func TestParseCompositeID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		wantParentID   int64
		wantResourceID int64
		wantErr        bool
	}{
		{
			name:           "valid composite ID",
			id:             "123:456",
			wantParentID:   123,
			wantResourceID: 456,
			wantErr:        false,
		},
		{
			name:    "missing colon",
			id:      "123456",
			wantErr: true,
		},
		{
			name:    "too many parts",
			id:      "123:456:789",
			wantErr: true,
		},
		{
			name:    "invalid parent ID",
			id:      "abc:456",
			wantErr: true,
		},
		{
			name:    "invalid resource ID",
			id:      "123:xyz",
			wantErr: true,
		},
		{
			name:    "zero IDs",
			id:      "0:0",
			wantErr: true,
		},
		{
			name:    "zero resource ID reads back as gone",
			id:      "123:0",
			wantErr: true,
		},
		{
			name:    "negative parent ID",
			id:      "-7:456",
			wantErr: true,
		},
		{
			// upapi.PrimaryKey is an int and the provider ships 386/arm builds, so an
			// oversized id would truncate and address a different record.
			name:    "parent ID beyond int64",
			id:      "99999999999999999999:456",
			wantErr: true,
		},
		{
			name:    "resource ID beyond int64",
			id:      "123:99999999999999999999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parentID, resourceID, err := ParseCompositeID(tt.id, "statuspage_id")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCompositeID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if parentID != tt.wantParentID {
					t.Errorf("ParseCompositeID() parentID = %v, want %v", parentID, tt.wantParentID)
				}
				if resourceID != tt.wantResourceID {
					t.Errorf("ParseCompositeID() resourceID = %v, want %v", resourceID, tt.wantResourceID)
				}
			}
		})
	}
}
