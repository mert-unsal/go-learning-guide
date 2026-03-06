package backtracking

import "testing"

func TestSubsetsWithDup(t *testing.T) {
	got := SubsetsWithDup([]int{1, 2, 2})
	if len(got) != 6 {
		t.Errorf("SubsetsWithDup([1,2,2]) returned %d subsets, want 6", len(got))
	}
}
