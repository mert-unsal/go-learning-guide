package linked_list

import (
	"reflect"
	"testing"
)

func TestReverseKGroup(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		k     int
		want  []int
	}{
		{"k=2", []int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{"k=3", []int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
		{"k=1", []int{1, 2, 3}, 1, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(ReverseKGroup(newList(tt.input), tt.k))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseKGroup(%v, %d) = %v, want %v", tt.input, tt.k, got, tt.want)
			}
		})
	}
}
