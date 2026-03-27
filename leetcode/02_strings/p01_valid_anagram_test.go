package strings_problems

import "testing"

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		name string
		s, t string
		want bool
	}{
		{"anagram", "anagram", "nagaram", true},
		{"not anagram", "rat", "car", false},
		{"empty", "", "", true},
		{"different lengths", "ab", "a", false},
		{"same chars diff count", "aa", "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAnagram(tt.s, tt.t)
			if got != tt.want {
				t.Errorf("IsAnagram(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
		})
	}
}
