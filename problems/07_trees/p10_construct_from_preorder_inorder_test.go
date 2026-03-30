package trees

import "testing"

func TestBuildTree(t *testing.T) {
	tests := []struct {
		name     string
		preorder []int
		inorder  []int
	}{
		{"basic", []int{3, 9, 20, 15, 7}, []int{9, 3, 15, 20, 7}},
		{"single", []int{1}, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildTree(tt.preorder, tt.inorder)
			_ = got // verify no panic
		})
	}
}
