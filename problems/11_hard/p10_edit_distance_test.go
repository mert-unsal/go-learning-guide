package hard

import "testing"

func TestMinDistance(t *testing.T) {
	tests := []struct {
		name         string
		word1, word2 string
		want         int
	}{
		{"horse to ros", "horse", "ros", 3},
		{"intention to execution", "intention", "execution", 5},
		{"empty to word", "", "abc", 3},
		{"word to empty", "abc", "", 3},
		{"same", "abc", "abc", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinDistance(tt.word1, tt.word2)
			if got != tt.want {
				t.Errorf("MinDistance(%q, %q) = %d, want %d", tt.word1, tt.word2, got, tt.want)
			}
		})
	}
}
