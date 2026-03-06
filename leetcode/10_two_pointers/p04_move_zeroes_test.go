package two_pointers

import (
	"reflect"
	"testing"
)

func TestMoveZeroes(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"basic", []int{0, 1, 0, 3, 12}, []int{1, 3, 12, 0, 0}},
		{"no zeros", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all zeros", []int{0, 0, 0}, []int{0, 0, 0}},
		{"single zero", []int{0}, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MoveZeroes(tt.input)
			if !reflect.DeepEqual(tt.input, tt.want) {
				t.Errorf("MoveZeroes result = %v, want %v", tt.input, tt.want)
			}
		})
	}
}
