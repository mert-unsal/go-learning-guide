package arrays

import "testing"

func TestMaximumToys(t *testing.T) {
	if got := MaximumToys([]int{1, 12, 5, 111, 200, 1000, 10}, 50); got != 4 {
		t.Errorf("got %d want 4", got)
	}
	if got := MaximumToys([]int{100, 200}, 50); got != 0 {
		t.Errorf("got %d want 0", got)
	}
}
