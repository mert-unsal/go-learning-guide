// Package strings_problems contains LeetCode string problems with explanations.
// Topics: frequency maps, sliding window, two-pointer on strings, Unicode.
//
// Note: this package is named strings_problems to avoid shadowing Go's
// built-in "strings" standard library package.
package strings_problems

import "strings"

// ============================================================
// PROBLEM 1: Valid Anagram (LeetCode #242) — EASY
// ============================================================
// Given two strings s and t, return true if t is an anagram of s.
// An anagram uses the same characters with the same frequencies.
//
// Example: s="anagram", t="nagaram" → true
// Example: s="rat",     t="car"     → false
//
// Approach: count character frequencies in s, then subtract for t.
// If all counts reach zero, they're anagrams.

// IsAnagram returns true if s and t are anagrams.
// Time: O(n)  Space: O(1) — only 26 lowercase letters
func IsAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false // different lengths can't be anagrams
	}
	var freq [26]int // index 0 = 'a', index 25 = 'z'

	for i := 0; i < len(s); i++ {
		freq[s[i]-'a']++ // increment for s
		freq[t[i]-'a']-- // decrement for t
	}

	// If s and t are anagrams, every frequency is exactly 0
	for _, count := range freq {
		if count != 0 {
			return false
		}
	}
	return true
}

// ============================================================
// PROBLEM 2: Longest Substring Without Repeating Characters (LeetCode #3) — MEDIUM
// ============================================================
// Find the length of the longest substring with all unique characters.
//
// Example: s="abcabcbb" → 3 ("abc")
// Example: s="pwwkew"   → 3 ("wke")
//
// Approach: sliding window.
// - right pointer expands the window by one character each step
// - if the new character is already in the window, shrink from left
//   until the duplicate is removed
// - Use a map: char → last-seen index (allows O(1) jump of left pointer)

// LengthOfLongestSubstring returns the length of the longest unique-char substring.
// Time: O(n)  Space: O(min(n, alphabet_size))
func LengthOfLongestSubstring(s string) int {
	// lastSeen maps each character to its most recent index
	lastSeen := make(map[byte]int)
	maxLen := 0
	left := 0

	for right := 0; right < len(s); right++ {
		ch := s[right]

		// If ch was seen and its last position is inside our window,
		// jump left pointer past that occurrence
		if idx, ok := lastSeen[ch]; ok && idx >= left {
			left = idx + 1
		}

		lastSeen[ch] = right // update last-seen position

		if right-left+1 > maxLen {
			maxLen = right - left + 1
		}
	}
	return maxLen
}

// ============================================================
// PROBLEM 3: Valid Palindrome (LeetCode #125) — EASY
// ============================================================
// A phrase is a palindrome if, after keeping only alphanumeric characters
// and converting to lowercase, it reads the same forwards and backwards.
//
// Example: "A man, a plan, a canal: Panama" → true
// Example: "race a car" → false
//
// Approach: two-pointer converging from both ends, skip non-alphanumeric.

// IsPalindrome returns true if s is a valid palindrome (ignoring case/non-alnum).
// Time: O(n)  Space: O(1)
func IsPalindrome(s string) bool {
	left, right := 0, len(s)-1

	for left < right {
		// Skip non-alphanumeric characters from the left
		for left < right && !isAlphanumeric(s[left]) {
			left++
		}
		// Skip non-alphanumeric characters from the right
		for left < right && !isAlphanumeric(s[right]) {
			right--
		}
		// Compare characters (case-insensitive)
		if toLower(s[left]) != toLower(s[right]) {
			return false
		}
		left++
		right--
	}
	return true
}

// isAlphanumeric returns true if b is a letter or digit.
func isAlphanumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9')
}

// toLower converts an ASCII byte to lowercase (no-op if already lowercase/digit).
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}

// ============================================================
// PROBLEM 4: Longest Common Prefix (LeetCode #14) — EASY
// ============================================================
// Find the longest common prefix string among an array of strings.
// Return "" if no common prefix exists.
//
// Example: ["flower","flow","flight"] → "fl"
//
// Approach: use the first string as a reference, shrink it character by character.

// LongestCommonPrefix returns the longest common prefix of strs.
// Time: O(S) where S = total characters  Space: O(1)
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]

	for _, s := range strs[1:] {
		// Trim prefix until s starts with it
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

// ============================================================
// PROBLEM 5: Reverse Words in a String (LeetCode #151) — MEDIUM
// ============================================================
// Given a string, reverse the order of words (trim extra spaces).
//
// Example: "  the sky is blue  " → "blue is sky the"
//
// Approach: use strings.Fields to split on any whitespace, then join reversed.

// ReverseWords reverses the word order in a string.
// Time: O(n)  Space: O(n)
func ReverseWords(s string) string {
	words := strings.Fields(s) // splits on any whitespace, trims
	left, right := 0, len(words)-1
	for left < right {
		words[left], words[right] = words[right], words[left]
		left++
		right--
	}
	return strings.Join(words, " ")
}
