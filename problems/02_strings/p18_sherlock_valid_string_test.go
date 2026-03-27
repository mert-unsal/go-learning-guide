package strings_problems

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"aabbcc", "YES"},
		{"aabbccc", "YES"},
		{"aabbccdd", "YES"},
		{"abcde", "YES"},
		{"aabbc", "YES"},
		{"aabbcd", "NO"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValid(tt.input); got != tt.want {
				t.Errorf("IsValid(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
