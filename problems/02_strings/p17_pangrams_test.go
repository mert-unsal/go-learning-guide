package strings_problems

import "testing"

func TestPangram(t *testing.T) {
	if got := Pangram("We promptly judged antique ivory buckles for the next prize"); got != "pangram" {
		t.Errorf("got %q want pangram", got)
	}
	if got := Pangram("The quick brown fox jumps over the lazy dog"); got != "pangram" {
		t.Errorf("got %q want pangram", got)
	}
	if got := Pangram("hello world"); got != "not pangram" {
		t.Errorf("got %q want not pangram", got)
	}
}
