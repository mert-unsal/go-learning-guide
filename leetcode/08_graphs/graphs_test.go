package graphs

import (
	"reflect"
	"testing"
)

// copyGrid deep-copies a [][]byte grid (NumIslands mutates it in-place)
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
		{
			"all water",
			[][]byte{{'0', '0'}, {'0', '0'}},
			0,
		},
		{
			"all land",
			[][]byte{{'1', '1'}, {'1', '1'}},
			1,
		},
		{
			"single cell",
			[][]byte{{'1'}},
			1,
		},
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

func TestCanFinish(t *testing.T) {
	tests := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		want          bool
	}{
		{"no prereqs", 2, [][]int{}, true},
		{"linear chain", 2, [][]int{{1, 0}}, true},
		{"simple cycle", 2, [][]int{{1, 0}, {0, 1}}, false},
		{"longer no cycle", 4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}, true},
		{"longer cycle", 3, [][]int{{0, 1}, {1, 2}, {2, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanFinish(tt.numCourses, tt.prerequisites)
			if got != tt.want {
				t.Errorf("CanFinish(%d, %v) = %v, want %v", tt.numCourses, tt.prerequisites, got, tt.want)
			}
		})
	}
}

func TestFloodFill(t *testing.T) {
	tests := []struct {
		name          string
		image         [][]int
		sr, sc, color int
		want          [][]int
	}{
		{
			"basic",
			[][]int{{1, 1, 1}, {1, 1, 0}, {1, 0, 1}},
			1, 1, 2,
			[][]int{{2, 2, 2}, {2, 2, 0}, {2, 0, 1}},
		},
		{
			"same color noop",
			[][]int{{0, 0, 0}, {0, 0, 0}},
			0, 0, 0,
			[][]int{{0, 0, 0}, {0, 0, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FloodFill(tt.image, tt.sr, tt.sc, tt.color)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FloodFill = %v, want %v", got, tt.want)
			}
		})
	}
}
