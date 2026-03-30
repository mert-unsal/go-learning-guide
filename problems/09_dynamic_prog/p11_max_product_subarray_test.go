package dynamic_prog

import "testing"

func TestMaxProduct(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{2, 3, -2, 4}, 6},
		{"negative", []int{-2, 0, -1}, 0},
		{"single negative", []int{-2}, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxProduct(tt.nums)
			if got != tt.want {
				t.Errorf("MaxProduct(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
