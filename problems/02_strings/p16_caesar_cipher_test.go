package strings_problems

import "testing"

func TestCaesarCipher(t *testing.T) {
	if got := CaesarCipher("middle-Outz", 2); got != "okffng-Qwvb" {
		t.Errorf("got %q want okffng-Qwvb", got)
	}
	if got := CaesarCipher("abc", 3); got != "def" {
		t.Errorf("got %q want def", got)
	}
	if got := CaesarCipher("xyz", 3); got != "abc" {
		t.Errorf("got %q want abc", got)
	}
}
