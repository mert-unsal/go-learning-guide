package trees

import "testing"

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
