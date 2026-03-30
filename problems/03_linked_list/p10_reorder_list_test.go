package linked_list

import (
	"reflect"
	"testing"
)

func TestReorderList(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want []int
	}{
		{"four elements", []int{1, 2, 3, 4}, []int{1, 4, 2, 3}},
		{"five elements", []int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}},
		{"single", []int{1}, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := newList(tt.vals)
			ReorderList(head)
			got := toSlice(head)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReorderList(%v) = %v, want %v", tt.vals, got, tt.want)
			}
		})
	}
}
