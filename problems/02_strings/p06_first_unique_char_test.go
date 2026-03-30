package strings_problems

import "testing"

func TestFirstUniqChar(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"first char unique", "leetcode", 0},
		{"middle unique", "loveleetcode", 2},
		{"no unique", "aabb", -1},
		{"single char", "z", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirstUniqChar(tt.s)
			if got != tt.want {
				t.Errorf("FirstUniqChar(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
