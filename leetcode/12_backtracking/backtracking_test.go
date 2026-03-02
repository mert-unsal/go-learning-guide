package backtracking

import (
	"reflect"
	"sort"
	"testing"
)

func TestSubsets(t *testing.T) {
	got := Subsets([]int{1, 2, 3})
	if len(got) != 8 {
		t.Errorf("Subsets([1,2,3]) returned %d subsets, want 8", len(got))
	}
}

func TestPermute(t *testing.T) {
	got := Permute([]int{1, 2, 3})
	if len(got) != 6 {
		t.Errorf("Permute([1,2,3]) returned %d permutations, want 6", len(got))
	}
}

func TestCombinationSum(t *testing.T) {
	got := CombinationSum([]int{2, 3, 6, 7}, 7)
	// Sort for comparison
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

func TestExist(t *testing.T) {
	board := [][]byte{
		{'A', 'B', 'C', 'E'},
		{'S', 'F', 'C', 'S'},
		{'A', 'D', 'E', 'E'},
	}
	if !Exist(board, "ABCCED") {
		t.Error("Exist(board, ABCCED) = false, want true")
	}
}

func TestLetterCombinations(t *testing.T) {
	got := LetterCombinations("23")
	want := []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LetterCombinations(23) = %v, want %v", got, want)
	}
}

func TestPartition(t *testing.T) {
	got := Partition("aab")
	if len(got) != 2 {
		t.Errorf("Partition(aab) returned %d partitions, want 2", len(got))
	}
}

func TestSubsetsWithDup(t *testing.T) {
	got := SubsetsWithDup([]int{1, 2, 2})
	if len(got) != 6 {
		t.Errorf("SubsetsWithDup([1,2,2]) returned %d subsets, want 6", len(got))
	}
}
