package graphs

import "testing"

func TestMaxAreaOfIsland(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"basic", [][]int{{0, 0, 1, 0, 0}, {0, 0, 0, 0, 0}, {0, 1, 1, 0, 0}, {0, 0, 0, 1, 1}}, 2},
		{"no island", [][]int{{0, 0, 0}, {0, 0, 0}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxAreaOfIsland(tt.grid)
			if got != tt.want {
				t.Errorf("MaxAreaOfIsland(%v) = %d, want %d", tt.grid, got, tt.want)
			}
		})
	}
}
