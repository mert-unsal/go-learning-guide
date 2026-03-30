package dynamic_prog

import "testing"

func TestIsInterleave(t *testing.T) {
	tests := []struct {
		name         string
		s1, s2, s3   string
		want         bool
	}{
		{"basic true", "aabcc", "dbbca", "aadbbcbcac", true},
		{"basic false", "aabcc", "dbbca", "aadbbbaccc", false},
		{"empty all", "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInterleave(tt.s1, tt.s2, tt.s3)
			if got != tt.want {
				t.Errorf("IsInterleave(%q, %q, %q) = %v, want %v",
					tt.s1, tt.s2, tt.s3, got, tt.want)
			}
		})
	}
}
