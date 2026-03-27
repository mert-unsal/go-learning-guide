package graphs

import (
	"reflect"
	"testing"
)

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
