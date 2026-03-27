package arrays

import "testing"

func TestArrayManipulation(t *testing.T) {
	if got := ArrayManipulation(5, [][]int{{1, 2, 100}, {2, 5, 100}, {3, 4, 100}}); got != 200 {
		t.Errorf("got %d want 200", got)
	}
	if got := ArrayManipulation(10, [][]int{{1, 5, 3}, {4, 8, 7}, {6, 9, 1}}); got != 10 {
		t.Errorf("got %d want 10", got)
	}
}
