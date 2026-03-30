package strings_problems

import (
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		strs []string
	}{
		{"basic", []string{"hello", "world"}},
		{"empty strings", []string{"", ""}},
		{"with delimiter", []string{"contains#delimiter"}},
		{"single", []string{"abc"}},
		{"empty list", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode(tt.strs)
			got := Decode(encoded)
			if !reflect.DeepEqual(got, tt.strs) {
				t.Errorf("Decode(Encode(%v)) = %v, want %v", tt.strs, got, tt.strs)
			}
		})
	}
}
