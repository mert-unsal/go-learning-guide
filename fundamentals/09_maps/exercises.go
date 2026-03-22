package maps

import "strings"

// ============================================================
// EXERCISES — 09 Maps
// ============================================================

// Exercise 1:
// CharFrequency counts the frequency of each character in a string.
//
// LESSON: `range` over a string gives (byte index, rune) — not (index, byte).
// Runes handle Unicode correctly. "hello" → map[rune]int{104:1, 101:1, 108:2, 111:1}
// (104='h', 101='e', 108='l', 111='o')
func CharFrequency(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, ch := range s { // ch is a rune (Unicode code point)
		freq[ch]++ // map zero-value is 0, so ++ works on first access
	}
	return freq
}

// Exercise 2:
// GroupByFirstChar groups words by their first character.
//
// LESSON: Map of slices — the value is itself a slice. You must append to it.
// The zero value of a slice is nil, and append(nil, x) works fine.
func GroupByFirstChar(words []string) map[byte][]string {
	groups := make(map[byte][]string)
	for _, w := range words {
		if len(w) > 0 {
			first := w[0]                            // w[0] is a byte — fine for ASCII first char
			groups[first] = append(groups[first], w) // append to nil is safe
		}
	}
	return groups
}

// Exercise 3:
// TopTwoFrequent returns the two most frequent elements.
//
// LESSON: Two-pass approach: first build a frequency map, then scan for top-2.
// Alternative: use a heap (see patterns/). For exactly top-2, scanning is simpler.
func TopTwoFrequent(nums []int) []int {
	freq := make(map[int]int)
	for _, n := range nums {
		freq[n]++
	}

	first, second := 0, 0
	firstFreq, secondFreq := 0, 0

	for num, count := range freq {
		if count > firstFreq {
			second, secondFreq = first, firstFreq
			first, firstFreq = num, count
		} else if count > secondFreq {
			second, secondFreq = num, count
		}
	}

	if secondFreq == 0 {
		return []int{first}
	}
	return []int{first, second}
}

// Exercise 4:
// IsAnagram checks if two strings are anagrams.
//
// LESSON: Classic map technique — count up for s, count down for t.
// If every count returns to 0, they use the same characters the same number of times.
// O(n) time, O(k) space where k = number of unique characters.
func IsAnagram(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	counts := make(map[rune]int)
	for _, ch := range s {
		counts[ch]++
	}
	for _, ch := range t {
		counts[ch]--
		if counts[ch] < 0 {
			return false // t has a char not in s, or more of it
		}
	}
	return true
}

// Exercise 5:
// FirstDuplicate returns the first integer that appears more than once.
//
// LESSON: "Seen set" pattern — use map[T]bool (or map[T]struct{}) as a set.
// map[T]struct{} uses slightly less memory (empty struct = 0 bytes),
// but map[T]bool is more readable.
func FirstDuplicate(nums []int) int {
	seen := make(map[int]bool)
	for _, n := range nums {
		if seen[n] { // zero value of bool is false, so this is safe
			return n
		}
		seen[n] = true
	}
	return -1
}

// Exercise 6:
// WordCount counts occurrences of each word in a sentence.
//
// LESSON: strings.Fields splits on any whitespace (spaces, tabs, newlines).
// Prefer it over strings.Split(s, " ") which doesn't handle multiple spaces.
func WordCount(sentence string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(sentence) {
		counts[word]++
	}
	return counts
}
