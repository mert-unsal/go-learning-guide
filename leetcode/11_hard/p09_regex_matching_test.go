package hard

import "testing"

func TestIsMatch(t *testing.T) {
	tests := []struct {
		name string
		s, p string
		want bool
	}{
		{"no match", "aa", "a", false},
		{"star zero or more", "aa", "a*", true},
		{"dot star", "ab", ".*", true},
		{"complex", "aab", "c*a*b", true},
		{"exact", "mississippi", "mis*is*p*.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMatch(tt.s, tt.p)
			if got != tt.want {
				t.Errorf("IsMatch(%q, %q) = %v, want %v", tt.s, tt.p, got, tt.want)
			}
		})
	}
}
