package arrays

import "testing"

func TestEraseOverlapIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals [][]int
		want      int
	}{
		{"basic", [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 1},
		{"all overlap", [][]int{{1, 2}, {1, 2}, {1, 2}}, 2},
		{"none overlap", [][]int{{1, 2}, {2, 3}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EraseOverlapIntervals(tt.intervals)
			if got != tt.want {
				t.Errorf("EraseOverlapIntervals() = %d, want %d", got, tt.want)
			}
		})
	}
}
