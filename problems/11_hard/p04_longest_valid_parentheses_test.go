package hard

import "testing"

func TestLongestValidParentheses(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"mixed", ")()())", 4},
		{"open left", "(()", 2},
		{"empty", "", 0},
		{"all valid", "()()", 4},
		{"nested", "((()))", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestValidParentheses(tt.s)
			if got != tt.want {
				t.Errorf("LongestValidParentheses(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
