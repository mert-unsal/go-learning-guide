package strings_problems

import "testing"

func TestRepeatedString(t *testing.T) {
	if got := RepeatedString("aba", 10); got != 7 {
		t.Errorf("got %d want 7", got)
	}
	if got := RepeatedString("a", 1000000000000); got != 1000000000000 {
		t.Errorf("got %d want 1000000000000", got)
	}
}
