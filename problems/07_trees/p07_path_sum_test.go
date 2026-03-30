package trees

import "testing"

func TestHasPathSum(t *testing.T) {
	tests := []struct {
		name      string
		vals      []int
		targetSum int
		want      bool
	}{
		{"has path", []int{5, 4, 8, 11, 0, 13, 4, 7, 2, 0, 0, 0, 1}, 22, true},
		{"no path", []int{1, 2, 3}, 5, false},
		{"empty tree", []int{}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasPathSum(newTree(tt.vals), tt.targetSum)
			if got != tt.want {
				t.Errorf("HasPathSum = %v, want %v", got, tt.want)
			}
		})
	}
}
