package strings_problems

import "testing"

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		name string
		s, t string
		want bool
	}{
		{"anagram", "anagram", "nagaram", true},
		{"not anagram", "rat", "car", false},
		{"empty", "", "", true},
		{"different lengths", "ab", "a", false},
		{"same chars diff count", "aa", "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAnagram(tt.s, tt.t)
			if got != tt.want {
				t.Errorf("IsAnagram(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
		})
	}
}

func TestLengthOfLongestSubstring(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"abc repeat", "abcabcbb", 3},
		{"all same", "bbbbb", 1},
		{"pwwkew", "pwwkew", 3},
		{"empty", "", 0},
		{"single", "a", 1},
		{"all unique", "abcdef", 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LengthOfLongestSubstring(tt.s)
			if got != tt.want {
				t.Errorf("LengthOfLongestSubstring(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

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

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want string
	}{
		{"flower", []string{"flower", "flow", "flight"}, "fl"},
		{"no common", []string{"dog", "racecar", "car"}, ""},
		{"empty slice", []string{}, ""},
		{"single", []string{"hello"}, "hello"},
		{"all same", []string{"abc", "abc", "abc"}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestCommonPrefix(tt.strs)
			if got != tt.want {
				t.Errorf("LongestCommonPrefix(%v) = %q, want %q", tt.strs, got, tt.want)
			}
		})
	}
}

func TestReverseWords(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"basic", "the sky is blue", "blue is sky the"},
		{"leading trailing spaces", "  hello world  ", "world hello"},
		{"multiple spaces", "a good   example", "example good a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReverseWords(tt.s)
			if got != tt.want {
				t.Errorf("ReverseWords(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}
