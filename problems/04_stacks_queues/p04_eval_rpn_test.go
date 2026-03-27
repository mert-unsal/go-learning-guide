package stacks_queues

import "testing"

func TestEvalRPN(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		want   int
	}{
		{"addition multiply", []string{"2", "1", "+", "3", "*"}, 9},
		{"division addition", []string{"4", "13", "5", "/", "+"}, 6},
		{"complex", []string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"}, 22},
		{"negative", []string{"3", "-4", "+"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvalRPN(tt.tokens)
			if got != tt.want {
				t.Errorf("EvalRPN(%v) = %d, want %d", tt.tokens, got, tt.want)
			}
		})
	}
}
