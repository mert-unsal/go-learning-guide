package trees

import "testing"

func TestLowestCommonAncestorBST(t *testing.T) {
	tests := []struct {
		name   string
		vals   []int
		pVal   int
		qVal   int
		expect int
	}{
		{"root is LCA", []int{6, 2, 8, 0, 4, 7, 9}, 2, 8, 6},
		{"subtree LCA", []int{6, 2, 8, 0, 4, 7, 9}, 2, 4, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newTree(tt.vals)
			_ = LowestCommonAncestorBST(root, root, root) // basic call test
		})
	}
}
