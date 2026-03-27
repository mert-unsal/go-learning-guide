package linked_list

import (
	"reflect"
	"testing"
)

func TestMergeTwoLists(t *testing.T) {
	tests := []struct {
		name  string
		list1 []int
		list2 []int
		want  []int
	}{
		{"normal", []int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
		{"both empty", []int{}, []int{}, nil},
		{"one empty", []int{}, []int{0}, []int{0}},
		{"first longer", []int{1, 3, 5, 7}, []int{2, 4}, []int{1, 2, 3, 4, 5, 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(MergeTwoLists(newList(tt.list1), newList(tt.list2)))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeTwoLists(%v, %v) = %v, want %v", tt.list1, tt.list2, got, tt.want)
			}
		})
	}
}
