package sliding_window

import "testing"

func TestMinWindow(t *testing.T) {
	tests := []struct {
		name string
		s, t string
		want string
	}{
		{"classic", "ADOBECODEBANC", "ABC", "BANC"},
		{"same string", "a", "a", "a"},
		{"no match", "a", "aa", ""},
		{"empty t", "a", "", ""},
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
