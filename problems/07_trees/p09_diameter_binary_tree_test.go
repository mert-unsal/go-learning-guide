package trees

import "testing"

func TestDiameterOfBinaryTree(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want int
	}{
		{"basic", []int{1, 2, 3, 4, 5}, 3},
		{"single", []int{1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DiameterOfBinaryTree(newTree(tt.vals))
			if got != tt.want {
				t.Errorf("DiameterOfBinaryTree = %v, want %v", got, tt.want)
			}
		})
	}
}
