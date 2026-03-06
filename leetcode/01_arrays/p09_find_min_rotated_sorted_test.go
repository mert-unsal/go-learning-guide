package arrays

import "testing"

func TestFindMinRotated(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"rotated", []int{3, 4, 5, 1, 2}, 1},
		{"not rotated", []int{1, 2, 3, 4, 5}, 1},
		{"single", []int{5}, 5},
		{"two", []int{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMinRotated(tt.nums)
			if got != tt.want {
				t.Errorf("FindMinRotated(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
