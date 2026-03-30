package linked_list

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want bool
	}{
		{"palindrome even", []int{1, 2, 2, 1}, true},
		{"not palindrome", []int{1, 2}, false},
		{"single", []int{1}, true},
		{"palindrome odd", []int{1, 2, 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(newList(tt.vals))
			if got != tt.want {
				t.Errorf("IsPalindrome(%v) = %v, want %v", tt.vals, got, tt.want)
			}
		})
	}
}
