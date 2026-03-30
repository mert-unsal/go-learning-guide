package two_pointers

import (
	"reflect"
	"testing"
)

func TestSortedSquares(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"mixed negatives", []int{-4, -1, 0, 3, 10}, []int{0, 1, 9, 16, 100}},
		{"all negative", []int{-7, -3, -1}, []int{1, 9, 49}},
		{"all positive", []int{1, 2, 3}, []int{1, 4, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortedSquares(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortedSquares(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
