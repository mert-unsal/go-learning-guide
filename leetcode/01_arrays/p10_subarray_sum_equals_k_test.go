package arrays

import "testing"

func TestSubarraySum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"basic", []int{1, 1, 1}, 2, 2},
		{"single", []int{1, 2, 3}, 3, 2},
		{"negative", []int{1, -1, 1}, 1, 3},
		{"zero k", []int{0, 0, 0}, 0, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SubarraySum(tt.nums, tt.k)
			if got != tt.want {
				t.Errorf("SubarraySum(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}
