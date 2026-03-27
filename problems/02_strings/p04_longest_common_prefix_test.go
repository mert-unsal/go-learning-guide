package strings_problems

import "testing"

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want string
	}{
		{"flower", []string{"flower", "flow", "flight"}, "fl"},
		{"no common", []string{"dog", "racecar", "car"}, ""},
		{"empty slice", []string{}, ""},
		{"single", []string{"hello"}, "hello"},
		{"all same", []string{"abc", "abc", "abc"}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestCommonPrefix(tt.strs)
			if got != tt.want {
				t.Errorf("LongestCommonPrefix(%v) = %q, want %q", tt.strs, got, tt.want)
			}
		})
	}
}
