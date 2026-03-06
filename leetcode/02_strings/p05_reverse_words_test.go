package strings_problems

import "testing"

func TestReverseWords(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"basic", "the sky is blue", "blue is sky the"},
		{"leading trailing spaces", "  hello world  ", "world hello"},
		{"multiple spaces", "a good   example", "example good a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReverseWords(tt.s)
			if got != tt.want {
				t.Errorf("ReverseWords(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}
