package stacks_queues

import "testing"

func TestGenerateParenthesis(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int // expected count of combinations
	}{
		{"n=1", 1, 1},
		{"n=2", 2, 2},
		{"n=3", 3, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateParenthesis(tt.n)
			if len(got) != tt.want {
				t.Errorf("GenerateParenthesis(%d) returned %d results, want %d", tt.n, len(got), tt.want)
			}
		})
	}
}
