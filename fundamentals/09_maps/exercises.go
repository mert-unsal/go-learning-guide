package maps

// ============================================================
// EXERCISES — 09 Maps
// ============================================================

// Exercise 1:
// CharFrequency counts the frequency of each character in a string.
//
// LESSON: `range` over a string gives (byte index, rune) — not (index, byte).
// Runes handle Unicode correctly. "hello" → map[rune]int{'h':1, 'e':1, 'l':2, 'o':1}
func CharFrequency(s string) map[rune]int {
	// TODO: iterate runes, count each in a map
	panic("not implemented")
}

// Exercise 2:
// GroupByFirstChar groups words by their first character.
//
// LESSON: Map of slices — the value is itself a slice. You must append to it.
// The zero value of a slice is nil, and append(nil, x) works fine.
func GroupByFirstChar(words []string) map[byte][]string {
	// TODO: group words by w[0] (first byte — fine for ASCII)
	panic("not implemented")
}

// Exercise 3:
// TopTwoFrequent returns the two most frequent elements.
//
// LESSON: Two-pass approach: first build a frequency map, then scan for top-2.
// Alternative: use a heap (see patterns/). For exactly top-2, scanning is simpler.
func TopTwoFrequent(nums []int) []int {
	// TODO: build frequency map, then find top two by count
	panic("not implemented")
}

// Exercise 4:
// IsAnagram checks if two strings are anagrams.
//
// LESSON: Classic map technique — count up for s, count down for t.
// If every count returns to 0, they use the same characters the same number of times.
// O(n) time, O(k) space where k = number of unique characters.
func IsAnagram(s, t string) bool {
	// TODO: count chars in s, decrement for t, check all zero
	panic("not implemented")
}

// Exercise 5:
// FirstDuplicate returns the first integer that appears more than once.
//
// LESSON: "Seen set" pattern — use map[T]bool (or map[T]struct{}) as a set.
// map[T]struct{} uses slightly less memory (empty struct = 0 bytes),
// but map[T]bool is more readable.
func FirstDuplicate(nums []int) int {
	// TODO: use a seen set, return first repeated element (or -1)
	panic("not implemented")
}

// Exercise 6:
// WordCount counts occurrences of each word in a sentence.
//
// LESSON: strings.Fields splits on any whitespace (spaces, tabs, newlines).
// Prefer it over strings.Split(s, " ") which doesn't handle multiple spaces.
func WordCount(sentence string) map[string]int {
	// TODO: split sentence into words, count each
	panic("not implemented")
}
