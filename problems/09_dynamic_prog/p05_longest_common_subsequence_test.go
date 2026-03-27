package dynamic_prog

import "testing"

func TestLongestCommonSubsequence(t *testing.T) {
	tests := []struct {
		name         string
		text1, text2 string
		want         int
	}{
		{"classic", "abcde", "ace", 3},
		{"same string", "abc", "abc", 3},
		{"no common", "abc", "def", 0},
		{"one empty", "", "abc", 0},
		{"partial", "oxcpqrsvwf", "shmtulqrypy", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestCommonSubsequence(tt.text1, tt.text2)
			if got != tt.want {
				t.Errorf("LCS(%q, %q) = %d, want %d", tt.text1, tt.text2, got, tt.want)
			}
		})
	}
}
