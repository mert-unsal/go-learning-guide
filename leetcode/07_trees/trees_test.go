package trees

import (
	"reflect"
	"testing"
)

func TestInorderTraversal(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want []int
	}{
		// Tree: 1 -> right=2, 2 -> left=3  =>  inorder: [1, 3, 2]
		{"right-then-left", []int{1, 0, 2, 3}, []int{1, 3, 2}},
		{"empty", []int{}, nil},
		{"single", []int{1}, []int{1}},
		// Tree: 3 -> left=1, 1 -> right=2  =>  inorder: [1, 2, 3]
		{"left only", []int{3, 1, 0, 0, 2}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newTree(tt.vals)
			got := InorderTraversal(root)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InorderTraversal = %v, want %v", got, tt.want)
			}
			// Also test iterative gives same result
			got2 := InorderIterative(root)
			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("InorderIterative = %v, want %v", got2, tt.want)
			}
		})
	}
}

func TestMaxDepth(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want int
	}{
		{"depth 3", []int{3, 9, 20, 0, 0, 15, 7}, 3},
		{"single", []int{1}, 1},
		{"empty", []int{}, 0},
		{"left chain", []int{1, 2, 0, 3}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxDepth(newTree(tt.vals))
			if got != tt.want {
				t.Errorf("MaxDepth = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestLevelOrder(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want [][]int
	}{
		{"normal", []int{3, 9, 20, 0, 0, 15, 7}, [][]int{{3}, {9, 20}, {15, 7}}},
		{"single", []int{1}, [][]int{{1}}},
		{"empty", []int{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LevelOrder(newTree(tt.vals))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevelOrder = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	// Build tree: 3, 5, 1, 6, 2, 0, 8, 0, 0, 7, 4
	//           3
	//          / \
	//         5   1
	//        / \ / \
	//       6  2 0  8
	//         / \
	//        7   4
	root := newTree([]int{3, 5, 1, 6, 2, 0, 8, 0, 0, 7, 4})
	// Find nodes with Val=5 and Val=1
	p := root.Left  // node 5
	q := root.Right // node 1
	lca := LowestCommonAncestor(root, p, q)
	if lca == nil || lca.Val != 3 {
		t.Errorf("LCA(5,1) = %v, want 3", lca)
	}
	// LCA of 5 and 4
	p2 := root.Left             // node 5
	q2 := root.Left.Right.Right // node 4
	lca2 := LowestCommonAncestor(root, p2, q2)
	if lca2 == nil || lca2.Val != 5 {
		t.Errorf("LCA(5,4) = %v, want 5", lca2)
	}
}
