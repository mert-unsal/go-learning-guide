package dynamic_prog

import "testing"

func TestRobII(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{2, 3, 2}, 3},
		{"four houses", []int{1, 2, 3, 1}, 4},
		{"single", []int{1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RobII(tt.nums)
			if got != tt.want {
				t.Errorf("RobII(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
