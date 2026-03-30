package dynamic_prog

import "testing"

func TestWordBreak(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		wordDict []string
		want     bool
	}{
		{"basic", "leetcode", []string{"leet", "code"}, true},
		{"reuse", "applepenapple", []string{"apple", "pen"}, true},
		{"impossible", "catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WordBreak(tt.s, tt.wordDict)
			if got != tt.want {
				t.Errorf("WordBreak(%q, %v) = %v, want %v", tt.s, tt.wordDict, got, tt.want)
			}
		})
	}
}
