package trees

import "testing"

func TestIsSymmetric(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want bool
	}{
		{"symmetric", []int{1, 2, 2, 3, 4, 4, 3}, true},
		{"not symmetric", []int{1, 2, 2, 0, 3, 0, 3}, false},
		{"single", []int{1}, true},
		{"empty", []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSymmetric(newTree(tt.vals))
			if got != tt.want {
				t.Errorf("IsSymmetric = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestLowestCommonAncestor(t *testing.T) {
	root := newTree([]int{3, 5, 1, 6, 2, 0, 8, 0, 0, 7, 4})
	p := root.Left
	q := root.Right
	lca := LowestCommonAncestor(root, p, q)
	if lca == nil || lca.Val != 3 {
		t.Errorf("LCA(5,1) = %v, want 3", lca)
	}
	p2 := root.Left
	q2 := root.Left.Right.Right
	lca2 := LowestCommonAncestor(root, p2, q2)
	if lca2 == nil || lca2.Val != 5 {
		t.Errorf("LCA(5,4) = %v, want 5", lca2)
	}
}
