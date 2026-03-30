package graphs

import "testing"

func TestPacificAtlantic(t *testing.T) {
	tests := []struct {
		name    string
		heights [][]int
	}{
		{"basic", [][]int{{1, 2, 2, 3, 5}, {3, 2, 3, 4, 4}, {2, 4, 5, 3, 1}, {6, 7, 1, 4, 5}, {5, 1, 1, 2, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: implement test validation
			_ = PacificAtlantic(tt.heights)
			t.Skip("not implemented")
		})
	}
}
