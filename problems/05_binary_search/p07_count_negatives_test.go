package binary_search

import "testing"

func TestCountNegatives(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"basic", [][]int{{4, 3, 2, -1}, {3, 2, 1, -1}, {1, 1, -1, -2}, {-1, -1, -2, -3}}, 8},
		{"all positive", [][]int{{3, 2}, {1, 0}}, 0},
		{"all negative", [][]int{{-1}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountNegatives(tt.grid)
			if got != tt.want {
				t.Errorf("CountNegatives = %v, want %v", got, tt.want)
			}
		})
	}
}
