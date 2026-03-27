package strings_problems

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"classic", "A man, a plan, a canal: Panama", true},
		{"not palindrome", "race a car", false},
		{"empty", " ", true},
		{"single char", "a", true},
		{"numbers", "0P", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(tt.s)
			if got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
