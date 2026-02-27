package maps
// ============================================================
// EXERCISES — 09 Maps
// ============================================================
// Exercise 1:
// Count the frequency of each character in a string.
// Return a map[rune]int.
// Example: "hello" → {h:1, e:1, l:2, o:1}
func CharFrequency(s string) map[rune]int {
// TODO: range over string gives (index, rune)
return nil
}
// Exercise 2:
// Given a slice of strings, group them by their first character.
// Return map[byte][]string.
// Example: ["ant","bat","bee","ape"] → {a:["ant","ape"], b:["bat","bee"]}
func GroupByFirstChar(words []string) map[byte][]string {
// TODO: use map[byte][]string, append to each group
return nil
}
// Exercise 3:
// Return the two most frequent elements in the slice.
// If there's a tie, return either order.
// Example: [1,1,2,2,3] → [1,2] (both have freq 2)
func TopTwoFrequent(nums []int) []int {
// TODO: count frequencies, then find top 2
return nil
}
// Exercise 4:
// Check if two strings are anagrams of each other.
// An anagram uses the same characters the same number of times.
// Example: "listen","silent" → true   "hello","world" → false
func IsAnagram(s, t string) bool {
// TODO: count chars in s, subtract for t, check all zero
return false
}
// Exercise 5:
// Given a slice of integers, return the first integer that appears
// more than once. Return -1 if no duplicates exist.
// Example: [4,3,2,7,8,2,3,1] → 2
func FirstDuplicate(nums []int) int {
// TODO: use a map[int]bool as a "seen" set
return -1
}
// Exercise 6:
// Word count — given a sentence, return a map of word→count.
// Words are separated by spaces.
// Example: "go is go" → {go:2, is:1}
func WordCount(sentence string) map[string]int {
// TODO: split by space, count each word
return nil
}