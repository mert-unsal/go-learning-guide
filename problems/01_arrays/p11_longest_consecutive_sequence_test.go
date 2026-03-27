package arrays

import "testing"

func TestLongestConsecutive(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{100, 4, 200, 1, 3, 2}, 4},
		{"single", []int{1}, 1},
		{"empty", []int{}, 0},
		{"duplicates", []int{1, 2, 0, 1}, 3},
		{"all same", []int{0, 0, 0}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestConsecutive(tt.nums)
			if got != tt.want {
				t.Errorf("LongestConsecutive(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
