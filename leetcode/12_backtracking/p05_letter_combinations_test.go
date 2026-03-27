package backtracking

import (
	"reflect"
	"testing"
)

func TestLetterCombinations(t *testing.T) {
	got := LetterCombinations("23")
	want := []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LetterCombinations(23) = %v, want %v", got, want)
	}
}
