// Maps Patterns — standalone demonstration of the five essential
// map-based patterns used in coding interviews and LeetCode.
//
// Run: go run ./cmd/concepts/maps/03-patterns
package main

import "fmt"

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

// ============================================================
// COMMON MAP PATTERNS (Essential for LeetCode!)
// ============================================================

// Pattern 1: Frequency counter (character/number frequency)
// Use case: anagram checks, character counting, histogram building.
// The zero value of int is 0, so freq[ch]++ works without initialization.
func FrequencyCount(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++ // if key doesn't exist, starts at zero value 0, then +1
	}
	return freq
}

// Pattern 2: Set (use map[T]bool or map[T]struct{})
// map[T]struct{} uses ZERO extra memory per entry because struct{} has
// size 0. Use this when you only care about membership, not a value.
// The comma-ok idiom checks membership without ambiguity.
func UniqueElements(nums []int) []int {
	seen := make(map[int]struct{}) // struct{} uses NO memory (vs bool)
	result := []int{}
	for _, n := range nums {
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			result = append(result, n)
		}
	}
	return result
}

// Pattern 3: Grouping (group items by a computed key)
// append() on a nil slice is safe — it allocates automatically.
// This pattern is used for grouping anagrams, categorizing data, etc.
func GroupByLength(words []string) map[int][]string {
	groups := make(map[int][]string)
	for _, w := range words {
		n := len(w)
		groups[n] = append(groups[n], w) // append to nil slice works!
	}
	return groups
}

// Pattern 4: Two-pass map (e.g., Two Sum — LeetCode #1)
// Build a map of value→index as you iterate. For each element,
// check if its complement (target - current) already exists.
// This turns O(n²) brute force into O(n) single-pass.
func TwoSum(nums []int, target int) (int, int) {
	// Map: value → index
	seen := make(map[int]int)
	for i, n := range nums {
		complement := target - n
		if j, ok := seen[complement]; ok {
			return j, i
		}
		seen[n] = i
	}
	return -1, -1
}

// Pattern 5: Memoization (cache expensive results)
// Maps are ideal for memoization because lookup is O(1) average.
// This converts exponential-time recursive Fibonacci into O(n).
// In production, consider sync.Map or mutex-protected maps for
// concurrent memoization.
func fibMemo(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	if v, ok := memo[n]; ok {
		return v
	}
	result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
	memo[n] = result
	return result
}

func Fibonacci(n int) int {
	return fibMemo(n, make(map[int]int))
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Maps: Essential Patterns               %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Pattern 1: Frequency Counter ---
	fmt.Printf("%s▸ Pattern 1: Frequency Counter%s\n", cyan+bold, reset)
	freq := FrequencyCount("hello world")
	fmt.Printf("  FrequencyCount(\"hello world\"):\n")
	for ch, count := range freq {
		fmt.Printf("    '%s%c%s' → %s%d%s\n", magenta, ch, reset, magenta, count, reset)
	}
	fmt.Printf("  %s✔ Zero value of int is 0 — freq[ch]++ works without initialization%s\n\n", green, reset)

	// --- Pattern 2: Set ---
	fmt.Printf("%s▸ Pattern 2: Set (map[T]struct{})%s\n", cyan+bold, reset)
	unique := UniqueElements([]int{1, 2, 2, 3, 3, 3, 4})
	fmt.Printf("  UniqueElements([1,2,2,3,3,3,4]) = %s%v%s\n", magenta, unique, reset)
	fmt.Printf("  %s✔ struct{} is zero-size — map[T]struct{} uses no memory per entry vs map[T]bool%s\n\n", green, reset)

	// --- Pattern 3: Grouping ---
	fmt.Printf("%s▸ Pattern 3: Grouping by Computed Key%s\n", cyan+bold, reset)
	groups := GroupByLength([]string{"go", "is", "fun", "and", "fast"})
	fmt.Printf("  GroupByLength([\"go\",\"is\",\"fun\",\"and\",\"fast\"]):\n")
	for length, words := range groups {
		fmt.Printf("    len=%s%d%s → %s%v%s\n", magenta, length, reset, magenta, words, reset)
	}
	fmt.Printf("  %s✔ append() on a nil slice is safe — auto-allocates on first append%s\n\n", green, reset)

	// --- Pattern 4: Two Sum ---
	fmt.Printf("%s▸ Pattern 4: Two Sum (hash map lookup)%s\n", cyan+bold, reset)
	nums := []int{2, 7, 11, 15}
	target := 9
	i, j := TwoSum(nums, target)
	fmt.Printf("  TwoSum(%v, %d) = indices [%s%d%s, %s%d%s]\n", nums, target, magenta, i, reset, magenta, j, reset)
	fmt.Printf("  %s✔ O(n) single-pass: for each element, check if complement exists in map%s\n", green, reset)
	fmt.Printf("  %s✔ Turns O(n²) brute force into O(n) with O(n) space%s\n\n", green, reset)

	// --- Pattern 5: Memoization ---
	fmt.Printf("%s▸ Pattern 5: Memoization (Fibonacci)%s\n", cyan+bold, reset)
	result := Fibonacci(10)
	fmt.Printf("  Fibonacci(10) = %s%d%s\n", magenta, result, reset)
	fmt.Printf("  %s✔ Map caches results: O(2^n) recursive → O(n) memoized%s\n", green, reset)
	fmt.Printf("  %s⚠ For concurrent memoization, use sync.Map or mutex-protected map%s\n", yellow, reset)
}
