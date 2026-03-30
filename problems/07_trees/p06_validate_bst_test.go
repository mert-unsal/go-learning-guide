package trees

import "testing"

func TestIsValidBST(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want bool
	}{
		{"valid BST", []int{2, 1, 3}, true},
		{"invalid", []int{5, 1, 4, 0, 0, 3, 6}, false},
		{"single", []int{1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidBST(newTree(tt.vals))
			if got != tt.want {
				t.Errorf("IsValidBST = %v, want %v", got, tt.want)
			}
		})
	}
}
