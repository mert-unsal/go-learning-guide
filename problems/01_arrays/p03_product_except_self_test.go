package arrays

import (
	"reflect"
	"testing"
)

func TestProductExceptSelf(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{"basic", []int{1, 2, 3, 4}, []int{24, 12, 8, 6}},
		{"with zero", []int{1, 0, 3, 4}, []int{0, 12, 0, 0}},
		{"two elements", []int{3, 4}, []int{4, 3}},
		{"negative", []int{-1, 1, 0, -3, 3}, []int{0, 0, 9, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProductExceptSelf(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductExceptSelf(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
