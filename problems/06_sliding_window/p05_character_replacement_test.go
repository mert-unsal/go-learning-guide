package sliding_window

import "testing"

func TestCharacterReplacement(t *testing.T) {
	tests := []struct {
		name string
		s    string
		k    int
		want int
	}{
		{"basic", "ABAB", 2, 4},
		{"one replace", "AABABBA", 1, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CharacterReplacement(tt.s, tt.k)
			if got != tt.want {
				t.Errorf("CharacterReplacement(%q, %d) = %v, want %v", tt.s, tt.k, got, tt.want)
			}
		})
	}
}
