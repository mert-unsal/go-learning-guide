package hard

import "testing"

func TestAlienOrder(t *testing.T) {
	tests := []struct {
		name      string
		words     []string
		wantEmpty bool
		wantLen   int
	}{
		{"basic", []string{"wrt", "wrf", "er", "ett", "rftt"}, false, 5},
		{"invalid prefix", []string{"abc", "ab"}, true, 0},
		{"simple", []string{"z", "x"}, false, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AlienOrder(tt.words)
			if tt.wantEmpty && got != "" {
				t.Errorf("AlienOrder(%v) = %q, want empty string", tt.words, got)
			}
			if !tt.wantEmpty && len(got) != tt.wantLen {
				t.Errorf("AlienOrder(%v) = %q (len %d), want len %d", tt.words, got, len(got), tt.wantLen)
			}
		})
	}
}
