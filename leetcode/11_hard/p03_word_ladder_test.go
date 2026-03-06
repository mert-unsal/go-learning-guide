package hard

import "testing"

func TestLadderLength(t *testing.T) {
	tests := []struct {
		name      string
		beginWord string
		endWord   string
		wordList  []string
		want      int
	}{
		{"basic", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}, 5},
		{"no path", "hit", "cog", []string{"hot", "dot", "dog", "lot", "log"}, 0},
		{"direct", "a", "c", []string{"a", "b", "c"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LadderLength(tt.beginWord, tt.endWord, tt.wordList)
			if got != tt.want {
				t.Errorf("LadderLength(%q, %q) = %d, want %d", tt.beginWord, tt.endWord, got, tt.want)
			}
		})
	}
}
