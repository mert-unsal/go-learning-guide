package trees

import (
	"reflect"
	"testing"
)

func TestInorderTraversal(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want []int
	}{
		{"right-then-left", []int{1, 0, 2, 3}, []int{1, 3, 2}},
		{"empty", []int{}, nil},
		{"single", []int{1}, []int{1}},
		{"left only", []int{3, 1, 0, 0, 2}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newTree(tt.vals)
			got := InorderTraversal(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InorderTraversal = %v, want %v", got, tt.want)
			}
			got2 := InorderIterative(root)
			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("InorderIterative = %v, want %v", got2, tt.want)
			}
		})
	}
}
