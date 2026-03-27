package binary_search

import "testing"

func TestSearchMatrix(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		target int
		want   bool
	}{
		{"found", [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 3, true},
		{"not found", [][]int{{1, 3, 5, 7}, {10, 11, 16, 20}, {23, 30, 34, 60}}, 13, false},
		{"single", [][]int{{1}}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchMatrix(tt.matrix, tt.target)
			if got != tt.want {
				t.Errorf("SearchMatrix target=%d = %v, want %v", tt.target, got, tt.want)
			}
		})
	}
}
