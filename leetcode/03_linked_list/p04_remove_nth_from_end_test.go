package linked_list

import (
	"reflect"
	"testing"
)

func TestRemoveNthFromEnd(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
		want  []int
	}{
		{"remove second from end", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 3, 5}},
		{"remove only node", []int{1}, 1, nil},
		{"remove last", []int{1, 2}, 1, []int{1}},
		{"remove first", []int{1, 2, 3}, 3, []int{2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(RemoveNthFromEnd(newList(tt.input), tt.n))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveNthFromEnd(%v, %d) = %v, want %v", tt.input, tt.n, got, tt.want)
			}
		})
	}
}
