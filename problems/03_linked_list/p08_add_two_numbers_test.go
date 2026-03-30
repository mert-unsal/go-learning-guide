package linked_list

import (
	"reflect"
	"testing"
)

func TestAddTwoNumbers(t *testing.T) {
	tests := []struct {
		name string
		l1   []int
		l2   []int
		want []int
	}{
		{"basic", []int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}},
		{"zeros", []int{0}, []int{0}, []int{0}},
		{"carry", []int{9, 9, 9}, []int{1}, []int{0, 0, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(AddTwoNumbers(newList(tt.l1), newList(tt.l2)))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddTwoNumbers(%v, %v) = %v, want %v", tt.l1, tt.l2, got, tt.want)
			}
		})
	}
}
