// Copyright Red Hat

package helpers

import (
	"os"
	"testing"
)

func TestGetExampleHeader(t *testing.T) {
	tests := []struct {
		name string
		arg0 string
		want string
	}{
		{
			name: "oc",
			arg0: "oc",
			want: "oc cm",
		},
		{
			name: "kubectl",
			arg0: "kubectl",
			want: "kubectl cm",
		},
		{
			name: "not-defined",
			arg0: "applier",
			want: "applier",
		},
	}
	for _, tt := range tests {
		os.Args[0] = tt.arg0
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExampleHeader(); got != tt.want {
				t.Errorf("GetExampleHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
