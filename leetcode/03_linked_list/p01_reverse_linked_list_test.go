package linked_list

import (
	"reflect"
	"testing"
)

func TestReverseList(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"normal", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"two nodes", []int{1, 2}, []int{2, 1}},
		{"single", []int{1}, []int{1}},
		{"empty", []int{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(ReverseList(newList(tt.input)))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseList(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestReverseListRecursive(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := []int{5, 4, 3, 2, 1}
	got := toSlice(ReverseListRecursive(newList(input)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReverseListRecursive(%v) = %v, want %v", input, got, want)
	}
}
