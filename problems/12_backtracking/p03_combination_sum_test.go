package backtracking

import (
	"sort"
	"testing"
)

func TestCombinationSum(t *testing.T) {
	got := CombinationSum([]int{2, 3, 6, 7}, 7)
	for _, combo := range got {
		sort.Ints(combo)
	}
	want := [][]int{{2, 2, 3}, {7}}
	if len(got) != len(want) {
		t.Errorf("CombinationSum([2,3,6,7], 7) returned %d combinations, want %d", len(got), len(want))
	}
}

func TestCombinationSum2(t *testing.T) {
	got := CombinationSum2([]int{10, 1, 2, 7, 6, 1, 5}, 8)
	if len(got) != 4 {
		t.Errorf("CombinationSum2 returned %d combinations, want 4", len(got))
	}
}
