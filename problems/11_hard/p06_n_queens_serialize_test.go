package hard

import "testing"

func makeTree(val int, left, right *TreeNode) *TreeNode {
	return &TreeNode{Val: val, Left: left, Right: right}
}

func TestSolveNQueens(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		wantCount int
	}{
		{"n=1", 1, 1},
		{"n=4", 4, 2},
		{"n=5", 5, 10},
		{"n=6", 6, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolveNQueens(tt.n)
			if len(got) != tt.wantCount {
				t.Errorf("SolveNQueens(%d) returned %d solutions, want %d", tt.n, len(got), tt.wantCount)
			}
		})
	}
}

func TestSerializeDeserialize(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
	}{
		{"simple tree", makeTree(1, makeTree(2, nil, nil), makeTree(3, makeTree(4, nil, nil), makeTree(5, nil, nil)))},
		{"nil", nil},
		{"single node", makeTree(42, nil, nil)},
		{"left skewed", makeTree(1, makeTree(2, makeTree(3, nil, nil), nil), nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized := Serialize(tt.root)
			deserialized := Deserialize(serialized)
			if Serialize(tt.root) != Serialize(deserialized) {
				t.Errorf("Serialize/Deserialize round-trip failed for %q", tt.name)
			}
		})
	}
}
