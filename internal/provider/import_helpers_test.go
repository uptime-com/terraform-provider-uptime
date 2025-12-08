package provider

import (
	"testing"
)

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
			name:           "zero IDs",
			id:             "0:0",
			wantParentID:   0,
			wantResourceID: 0,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parentID, resourceID, err := ParseCompositeID(tt.id)
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
