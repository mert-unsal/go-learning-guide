package strings_problems

import "testing"

func TestLengthOfLongestSubstring(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"abc repeat", "abcabcbb", 3},
		{"all same", "bbbbb", 1},
		{"pwwkew", "pwwkew", 3},
		{"empty", "", 0},
		{"single", "a", 1},
		{"all unique", "abcdef", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LengthOfLongestSubstring(tt.s)
			if got != tt.want {
				t.Errorf("LengthOfLongestSubstring(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
