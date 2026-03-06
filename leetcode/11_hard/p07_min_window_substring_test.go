package hard

import "testing"

func TestMinWindow(t *testing.T) {
	tests := []struct {
		name string
		s, t string
		want string
	}{
		{"basic", "ADOBECODEBANC", "ABC", "BANC"},
		{"same", "a", "a", "a"},
		{"not found", "a", "aa", ""},
		{"exact match", "ABC", "ABC", "ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinWindow(tt.s, tt.t)
			if got != tt.want {
				t.Errorf("MinWindow(%q, %q) = %q, want %q", tt.s, tt.t, got, tt.want)
			}
		})
	}
}
