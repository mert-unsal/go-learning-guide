package arrays

import (
	"reflect"
	"testing"
)

func TestRotate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"k=3", []int{1, 2, 3, 4, 5, 6, 7}, 3, []int{5, 6, 7, 1, 2, 3, 4}},
		{"k=2", []int{-1, -100, 3, 99}, 2, []int{3, 99, -1, -100}},
		{"k=0", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"k=len", []int{1, 2, 3}, 3, []int{1, 2, 3}},
		{"k>len", []int{1, 2, 3}, 5, []int{2, 3, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Rotate(tt.nums, tt.k)
			if !reflect.DeepEqual(tt.nums, tt.want) {
				t.Errorf("Rotate() = %v, want %v", tt.nums, tt.want)
			}
		})
	}
}
