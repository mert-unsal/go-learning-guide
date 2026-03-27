package binary_search

import "testing"

func TestFindMin(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"rotated 5", []int{3, 4, 5, 1, 2}, 1},
		{"rotated 7", []int{4, 5, 6, 7, 0, 1, 2}, 0},
		{"not rotated", []int{11, 13, 15, 17}, 11},
		{"single", []int{1}, 1},
		{"two elements rotated", []int{2, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMin(tt.nums)
			if got != tt.want {
				t.Errorf("FindMin(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
