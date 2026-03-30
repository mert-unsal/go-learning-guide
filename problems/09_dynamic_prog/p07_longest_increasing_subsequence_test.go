package dynamic_prog

import "testing"

func TestLengthOfLIS(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
		{"all increasing", []int{1, 2, 3, 4, 5}, 5},
		{"all decreasing", []int{5, 4, 3, 2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LengthOfLIS(tt.nums)
			if got != tt.want {
				t.Errorf("LengthOfLIS(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
