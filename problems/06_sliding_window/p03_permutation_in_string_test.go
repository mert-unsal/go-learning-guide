package sliding_window

import "testing"

func TestCheckInclusion(t *testing.T) {
	tests := []struct {
		name   string
		s1, s2 string
		want   bool
	}{
		{"contains perm", "ab", "eidbaooo", true},
		{"no perm", "ab", "eidboaoo", false},
		{"s1 longer", "abc", "ab", false},
		{"exact match", "abc", "cba", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckInclusion(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("CheckInclusion(%q, %q) = %v, want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}
