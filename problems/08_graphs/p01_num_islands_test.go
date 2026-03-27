package graphs

import "testing"

// copyGrid deep-copies a [][]byte grid (NumIslands mutates it in-place).
func copyGrid(grid [][]byte) [][]byte {
	cp := make([][]byte, len(grid))
	for i, row := range grid {
		cp[i] = make([]byte, len(row))
		copy(cp[i], row)
	}
	return cp
}

func TestNumIslands(t *testing.T) {
	tests := []struct {
		name string
		grid [][]byte
		want int
	}{
		{
			"two islands",
			[][]byte{
				{'1', '1', '0'},
				{'1', '1', '0'},
				{'0', '0', '1'},
			},
			2,
		},
		{
			"four islands",
			[][]byte{
				{'1', '0', '1'},
				{'0', '1', '0'},
				{'1', '0', '1'},
			},
			5,
		},
		{"all water", [][]byte{{'0', '0'}, {'0', '0'}}, 0},
		{"all land", [][]byte{{'1', '1'}, {'1', '1'}}, 1},
		{"single cell", [][]byte{{'1'}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NumIslands(copyGrid(tt.grid))
			if got != tt.want {
				t.Errorf("NumIslands = %d, want %d", got, tt.want)
			}
		})
	}
}
