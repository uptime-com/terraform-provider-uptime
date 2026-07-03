package provider

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/uptime-com/uptime-client-go/v2/pkg/upapi"
)

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "resource gone sentinel",
			err:  errResourceGone,
			want: true,
		},
		{
			name: "wrapped resource gone sentinel",
			err:  fmt.Errorf("read failed: %w", errResourceGone),
			want: true,
		},
		{
			name: "http 404",
			err:  &upapi.Error{Response: &http.Response{StatusCode: http.StatusNotFound}},
			want: true,
		},
		{
			name: "http 500",
			err:  &upapi.Error{Response: &http.Response{StatusCode: http.StatusInternalServerError}},
			want: false,
		},
		{
			name: "unrelated error",
			err:  errors.New("boom"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNotFoundError(tt.err); got != tt.want {
				t.Errorf("isNotFoundError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
