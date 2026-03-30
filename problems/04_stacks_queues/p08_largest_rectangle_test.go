package stacks_queues

import "testing"

func TestLargestRectangleArea(t *testing.T) {
	tests := []struct {
		name    string
		heights []int
		want    int
	}{
		{"basic", []int{2, 1, 5, 6, 2, 3}, 10},
		{"uniform", []int{2, 2, 2}, 6},
		{"single", []int{5}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LargestRectangleArea(tt.heights)
			if got != tt.want {
				t.Errorf("LargestRectangleArea(%v) = %v, want %v", tt.heights, got, tt.want)
			}
		})
	}
}
