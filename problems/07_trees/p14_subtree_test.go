package trees

import "testing"

func TestIsSubtree(t *testing.T) {
	tests := []struct {
		name    string
		root    []int
		subRoot []int
		want    bool
	}{
		{"is subtree", []int{3, 4, 5, 1, 2}, []int{4, 1, 2}, true},
		{"not subtree", []int{3, 4, 5, 1, 2, 0, 0, 0, 0, 0, 0}, []int{4, 1, 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSubtree(newTree(tt.root), newTree(tt.subRoot))
			if got != tt.want {
				t.Errorf("IsSubtree = %v, want %v", got, tt.want)
			}
		})
	}
}
