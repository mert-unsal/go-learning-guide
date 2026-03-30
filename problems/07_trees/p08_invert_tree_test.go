package trees

import "testing"

func TestInvertTree(t *testing.T) {
	tests := []struct {
		name string
		vals []int
	}{
		{"basic", []int{4, 2, 7, 1, 3, 6, 9}},
		{"single", []int{1}},
		{"empty", []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InvertTree(newTree(tt.vals))
			_ = got // verify no panic
		})
	}
}
