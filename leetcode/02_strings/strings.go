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

// ============================================================
// PROBLEM 6: First Unique Character in a String (LeetCode #387) — EASY
// ============================================================
// Find the first non-repeating character and return its index.
// Return -1 if none exists.
//
// Example: s="leetcode" → 0 ('l')
// Example: s="aabb"     → -1
//
// Approach: two-pass with frequency array.

// FirstUniqChar returns the index of the first non-repeating character.
// Time: O(n)  Space: O(1)
func FirstUniqChar(s string) int {
	var freq [26]int
	for _, ch := range s {
		freq[ch-'a']++
	}
	for i, ch := range s {
		if freq[ch-'a'] == 1 {
			return i
		}
	}
	return -1
}

// ============================================================
// PROBLEM 7: Roman to Integer (LeetCode #13) — EASY
// ============================================================
// Convert a Roman numeral string to an integer.
//
// Example: "III" → 3, "IV" → 4, "IX" → 9, "LVIII" → 58
//
// Key insight: if a smaller value appears BEFORE a larger value, subtract it.
// Otherwise add it.

// RomanToInt converts a Roman numeral string to an integer.
// Time: O(n)  Space: O(1)
func RomanToInt(s string) int {
	values := map[byte]int{
		'I': 1, 'V': 5, 'X': 10, 'L': 50,
		'C': 100, 'D': 500, 'M': 1000,
	}
	total := 0
	for i := 0; i < len(s); i++ {
		curr := values[s[i]]
		if i+1 < len(s) && curr < values[s[i+1]] {
			total -= curr // subtract (e.g. IV = 5-1)
		} else {
			total += curr
		}
	}
	return total
}

// ============================================================
// PROBLEM 8: Count and Say (LeetCode #38) — MEDIUM
// ============================================================
// The count-and-say sequence: describe the previous term.
// "1" → "11" (one 1) → "21" (two 1s) → "1211" → "111221" → ...
//
// Return the nth term.

// CountAndSay returns the nth term of the count-and-say sequence.
// Time: O(n * len of each term)  Space: O(n)
func CountAndSay(n int) string {
	result := "1"
	for i := 1; i < n; i++ {
		s := result
		result = ""
		j := 0
		for j < len(s) {
			ch := s[j]
			count := 1
			for j+count < len(s) && s[j+count] == ch {
				count++
			}
			// Append count then character
			result += string(rune('0'+count)) + string(ch)
			j += count
		}
	}
	return result
}

// ============================================================
// PROBLEM 9: Group Anagrams (LeetCode #49) — MEDIUM
// ============================================================
// Given an array of strings, group anagrams together.
//
// Example: ["eat","tea","tan","ate","nat","bat"]
//        → [["bat"],["nat","tan"],["ate","eat","tea"]]
//
// Key: sort each string → anagrams will have the same sorted key.

// GroupAnagrams groups strings that are anagrams of each other.
// Time: O(n * k log k) where k is max string length  Space: O(n*k)
func GroupAnagrams(strs []string) [][]string {
	groups := make(map[[26]int][]string)
	for _, s := range strs {
		var key [26]int
		for _, ch := range s {
			key[ch-'a']++
		}
		groups[key] = append(groups[key], s)
	}
	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	return result
}

// ============================================================
// PROBLEM 10: Encode and Decode Strings (LeetCode #271) — MEDIUM
// ============================================================
// Design an algorithm to encode a list of strings to a single string
// and decode it back. Must handle any characters including the delimiter.
//
// Approach: length-prefix encoding: "4#word4#next"
// Each string is prefixed by its length + '#'.

// Encode encodes a list of strings to a single string.
func Encode(strs []string) string {
	var sb strings.Builder
	for _, s := range strs {
		// Write: len(s) + '#' + s
		for _, d := range []byte(string(rune(len(s)))) {
			_ = d
		}
		sb.WriteString(strings.Repeat("", 0))
		length := len(s)
		// Build length prefix manually
		digits := ""
		if length == 0 {
			digits = "0"
		} else {
			tmp := length
			for tmp > 0 {
				digits = string(rune('0'+tmp%10)) + digits
				tmp /= 10
			}
		}
		sb.WriteString(digits)
		sb.WriteByte('#')
		sb.WriteString(s)
	}
	return sb.String()
}

// Decode decodes the encoded string back to a list of strings.
func Decode(s string) []string {
	var result []string
	i := 0
	for i < len(s) {
		// Read length digits until '#'
		j := i
		for s[j] != '#' {
			j++
		}
		length := 0
		for _, ch := range s[i:j] {
			length = length*10 + int(ch-'0')
		}
		i = j + 1
		result = append(result, s[i:i+length])
		i += length
	}
	return result
}
