package arrays

import "testing"

func TestAlmostSorted(t *testing.T) {
	if got := AlmostSorted([]int{2, 1}); got != "swap 1 2" {
		t.Errorf("[2,1] got %q want swap 1 2", got)
	}
	if got := AlmostSorted([]int{1, 5, 4, 3, 2, 6}); got != "reverse 2 5" {
		t.Errorf("[1,5,4,3,2,6] got %q want reverse 2 5", got)
	}
	if got := AlmostSorted([]int{1, 2, 3}); got != "yes" {
		t.Errorf("[1,2,3] got %q want yes", got)
	}
	if got := AlmostSorted([]int{3, 4, 1, 2}); got != "no" {
		t.Errorf("[3,4,1,2] got %q want no", got)
	}
	if got := AlmostSorted([]int{3, 1, 2}); got != "no" {
		t.Errorf("[3,1,2] got %q want no", got)
	}
}
