package maps
import (
"sort"
"strings"
)
// ============================================================
// SOLUTIONS — 09 Maps
// ============================================================
func CharFrequencySolution(s string) map[rune]int {
freq := make(map[rune]int)
for _, ch := range s { // range over string gives Unicode code points (runes)
freq[ch]++
}
return freq
}
func GroupByFirstCharSolution(words []string) map[byte][]string {
groups := make(map[byte][]string)
for _, word := range words {
if len(word) == 0 {
continue
}
first := word[0] // byte — works for ASCII first char
groups[first] = append(groups[first], word)
}
return groups
}
func TopTwoFrequentSolution(nums []int) []int {
freq := make(map[int]int)
for _, n := range nums {
freq[n]++
}
// collect unique numbers sorted by frequency descending
type pair struct{ val, count int }
pairs := make([]pair, 0, len(freq))
for v, c := range freq {
pairs = append(pairs, pair{v, c})
}
sort.Slice(pairs, func(i, j int) bool {
return pairs[i].count > pairs[j].count
})
result := make([]int, 0, 2)
for i := 0; i < 2 && i < len(pairs); i++ {
result = append(result, pairs[i].val)
}
return result
}
func IsAnagramSolution(s, t string) bool {
if len(s) != len(t) {
return false
}
count := make(map[rune]int)
for _, ch := range s {
count[ch]++
}
for _, ch := range t {
count[ch]--
if count[ch] < 0 {
return false // t has more of this char than s
}
}
return true
}
func FirstDuplicateSolution(nums []int) int {
seen := make(map[int]bool)
for _, n := range nums {
if seen[n] {
return n // first time we see it for the second time
}
seen[n] = true
}
return -1
}
func WordCountSolution(sentence string) map[string]int {
counts := make(map[string]int)
words := strings.Fields(sentence) // splits on any whitespace
for _, word := range words {
counts[word]++
}
return counts
}