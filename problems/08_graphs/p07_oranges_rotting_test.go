package graphs

import "testing"

func TestOrangesRotting(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"basic", [][]int{{2, 1, 1}, {1, 1, 0}, {0, 1, 1}}, 4},
		{"impossible", [][]int{{2, 1, 1}, {0, 1, 1}, {1, 0, 1}}, -1},
		{"already rotten", [][]int{{0, 2}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OrangesRotting(tt.grid)
			if got != tt.want {
				t.Errorf("OrangesRotting(%v) = %d, want %d", tt.grid, got, tt.want)
			}
		})
	}
}
