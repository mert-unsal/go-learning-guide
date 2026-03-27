package binary_search

import "testing"

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{"found", []int{-1, 0, 3, 5, 9, 12}, 9, 4},
		{"not found", []int{-1, 0, 3, 5, 9, 12}, 2, -1},
		{"empty", []int{}, 5, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(tt.nums, tt.target)
			if got != tt.want {
				t.Errorf("BinarySearch(%v, %d) = %d, want %d", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}
