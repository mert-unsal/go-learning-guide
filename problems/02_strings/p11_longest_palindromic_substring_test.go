package strings_problems

import "testing"

func TestLongestPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string // multiple valid answers
	}{
		{"babad", "babad", []string{"bab", "aba"}},
		{"cbbd", "cbbd", []string{"bb"}},
		{"single", "a", []string{"a"}},
		{"two same", "aa", []string{"aa"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestPalindrome(tt.s)
			valid := false
			for _, w := range tt.want {
				if got == w {
					valid = true
					break
				}
			}
			if !valid {
				t.Errorf("LongestPalindrome(%q) = %q, want one of %v", tt.s, got, tt.want)
			}
		})
	}
}
