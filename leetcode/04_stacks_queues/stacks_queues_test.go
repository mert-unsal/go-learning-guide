package stacks_queues

import (
	"reflect"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"simple pair", "()", true},
		{"multiple types", "()[]{}", true},
		{"wrong order", "(]", false},
		{"interleaved", "([)]", false},
		{"nested", "{[]}", true},
		{"empty", "", true},
		{"only open", "(((", false},
		{"only close", ")))", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValid(tt.s)
			if got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestMinStack(t *testing.T) {
	s := &MinStack{}
	s.Push(-2)
	s.Push(0)
	s.Push(-3)

	if got := s.GetMin(); got != -3 {
		t.Errorf("GetMin() = %d, want -3", got)
	}
	s.Pop()
	if got := s.Top(); got != 0 {
		t.Errorf("Top() = %d, want 0", got)
	}
	if got := s.GetMin(); got != -2 {
		t.Errorf("GetMin() after pop = %d, want -2", got)
	}
}

func TestDailyTemperatures(t *testing.T) {
	tests := []struct {
		name  string
		temps []int
		want  []int
	}{
		{"basic", []int{73, 74, 75, 71, 69, 72, 76, 73}, []int{1, 1, 4, 2, 1, 1, 0, 0}},
		{"all decreasing", []int{5, 4, 3, 2, 1}, []int{0, 0, 0, 0, 0}},
		{"all increasing", []int{1, 2, 3, 4, 5}, []int{1, 1, 1, 1, 0}},
		{"single", []int{30}, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DailyTemperatures(tt.temps)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DailyTemperatures(%v) = %v, want %v", tt.temps, got, tt.want)
			}
		})
	}
}

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
