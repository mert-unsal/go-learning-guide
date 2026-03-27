package strings_problems

import "testing"

func TestCountSubstrings(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"abc", "abc", 3},
		{"aaa", "aaa", 6},
		{"single", "a", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountSubstrings(tt.s)
			if got != tt.want {
				t.Errorf("CountSubstrings(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
