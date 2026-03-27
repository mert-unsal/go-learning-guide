package two_pointers

import (
	"reflect"
	"testing"
)

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"classic", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"all zeros", []int{0, 0, 0}, [][]int{{0, 0, 0}}},
		{"no solution", []int{1, 2, 3}, nil},
		{"with duplicates", []int{-2, 0, 0, 2, 2}, [][]int{{-2, 0, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ThreeSum(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ThreeSum(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
