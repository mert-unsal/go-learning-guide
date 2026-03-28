// Maps Patterns — standalone demonstration of the five essential
// map-based patterns used in coding interviews and LeetCode.
//
// Run: go run ./cmd/concepts/maps/03-patterns
package main

import "fmt"

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
	// Frequency counter
	freq := FrequencyCount("hello world")
	fmt.Println("Frequency:", freq)

	// Unique elements (set pattern)
	fmt.Println("Unique:", UniqueElements([]int{1, 2, 2, 3, 3, 3, 4}))

	// Grouping by computed key
	groups := GroupByLength([]string{"go", "is", "fun", "and", "fast"})
	fmt.Println("Groups by length:", groups)

	// Two Sum — single-pass hash map solution
	i, j := TwoSum([]int{2, 7, 11, 15}, 9)
	fmt.Println("TwoSum indices:", i, j) // 0, 1

	// Fibonacci with memoization
	fmt.Println("Fib(10):", Fibonacci(10)) // 55
}
