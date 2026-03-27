package backtracking

import "testing"

func TestSubsets(t *testing.T) {
	got := Subsets([]int{1, 2, 3})
	if len(got) != 8 {
		t.Errorf("Subsets([1,2,3]) returned %d subsets, want 8", len(got))
	}
}
