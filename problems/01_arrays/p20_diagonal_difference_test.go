package arrays

import "testing"

func TestDiagonalDifference(t *testing.T) {
	if got := DiagonalDifference([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}); got != 0 {
		t.Errorf("got %d want 0", got)
	}
	if got := DiagonalDifference([][]int{{11, 2, 4}, {4, 5, 6}, {10, 8, -12}}); got != 15 {
		t.Errorf("got %d want 15", got)
	}
}
