package two_pointers

import "testing"

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{1, 1, 2}, 2},
		{"multiple dups", []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}, 5},
		{"no dups", []int{1, 2, 3}, 3},
		{"single", []int{1}, 1},
		{"empty", []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDuplicates(tt.nums)
			if got != tt.want {
				t.Errorf("RemoveDuplicates(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
