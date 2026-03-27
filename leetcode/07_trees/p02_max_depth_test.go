package trees

import "testing"

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
