// Package sliding_window contains LeetCode sliding window problems.
// Topics: variable-size window, fixed-size window, character frequency maps.
package sliding_window

// ============================================================
// SLIDING WINDOW — Core Idea
// ============================================================
// A sliding window is a subarray/substring that moves through the input.
// Two pointers (left, right) define the window's boundaries.
//
// Template:
//   for right := 0; right < n; right++ {
//       add nums[right] to window state
//       for window is invalid {
//           remove nums[left] from window state
//           left++
//       }
//       update answer (window [left..right] is now valid)
//   }

// ============================================================
// PROBLEM 1: Maximum Average Subarray I (LeetCode #643) — EASY
// ============================================================
// Find the contiguous subarray of length k with the maximum average.
//
// Example: nums=[1,12,-5,-6,50,3], k=4 → 12.75  (subarray [12,-5,-6,50])
//
// Approach: fixed-size window. Slide a window of size k across the array.
// Track the running sum; update max sum whenever the window is full.

// FindMaxAverage returns the maximum average of any contiguous subarray of length k.
// Time: O(n)  Space: O(1)
func FindMaxAverage(nums []int, k int) float64 {
	// TODO: implement
	return 0
}

// ============================================================
// PROBLEM 2: Minimum Window Substring (LeetCode #76) — HARD
// ============================================================
// Given strings s and t, find the minimum window substring of s
// that contains all characters of t. Return "" if no such window exists.
//
// Example: s="ADOBECODEBANC", t="ABC" → "BANC"
//
// Approach: sliding window with frequency maps.
// - need: frequency of each char in t
// - have: frequency of each char in current window
// - formed: how many chars in t satisfy their required frequency
// Expand right, shrink left when window is valid (formed == len(unique t chars))

// MinWindow returns the minimum window substring containing all chars of t.
// Time: O(|s| + |t|)  Space: O(|s| + |t|)
func MinWindow(s string, t string) string {
	// TODO: implement
	return ""
}

// ============================================================
// PROBLEM 3: Permutation in String (LeetCode #567) — MEDIUM
// ============================================================
// Return true if s2 contains a permutation of s1.
// Equivalently: does any window of size len(s1) in s2 have the same
// character frequencies as s1?
//
// Approach: fixed-size sliding window (size = len(s1)).
// Compare character frequency arrays of window vs s1.

// CheckInclusion returns true if any permutation of s1 is a substring of s2.
// Time: O(|s1| + |s2|)  Space: O(1) — only 26 lowercase letters
func CheckInclusion(s1 string, s2 string) bool {
	// TODO: implement
	return false
}

// ============================================================
// PROBLEM 4: Fruit Into Baskets (LeetCode #904) — MEDIUM
// ============================================================
// You have two baskets, each holding one type of fruit.
// Find the max number of fruits you can pick from a contiguous subarray
// using at most 2 distinct fruit types.
//
// Essentially: find the length of the longest subarray with at most 2 distinct values.

// TotalFruit returns the max fruits pickable with at most 2 baskets.
// Time: O(n)  Space: O(1) — at most 3 entries in the map at once
func TotalFruit(fruits []int) int {
	// TODO: implement
	return 0
}

// ============================================================
// PROBLEM 5: Longest Repeating Character Replacement (LeetCode #424) — MEDIUM
// ============================================================
// You can replace at most k characters in a string. Find the length of
// the longest substring containing the same letter after replacements.
//
// Example: s="AABABBA", k=1 → 4
//
// Key insight: window is valid when (windowSize - maxFreq) <= k.
// We only care about the maximum frequency character in the window.
// We don't need to decrease maxFreq when shrinking — a smaller maxFreq
// only leads to smaller valid windows, which we don't care about.

// CharacterReplacement returns the longest valid substring length after k replacements.
// Time: O(n)  Space: O(1) — 26-char array
func CharacterReplacement(s string, k int) int {
	// TODO: implement
	return 0
}

// ============================================================
// PROBLEM 6: Maximum Points You Can Obtain from Cards (LeetCode #1423) — MEDIUM
// ============================================================
// From an array of card points, pick exactly k cards from either end.
// Maximize the total points.
//
// Example: cardPoints=[1,2,3,4,5,6,1], k=3 → 12  (pick 1,6,5)
//
// Key insight: total points = sum of all - minimum subarray of length (n-k).
// Find the minimum window of size n-k using a sliding window.

// MaxScore returns the maximum score by picking k cards from the ends.
// Time: O(n)  Space: O(1)
func MaxScore(cardPoints []int, k int) int {
	// TODO: implement
	return 0
}

// ============================================================
// PROBLEM 7: Minimum Size Subarray Sum (LeetCode #209) — MEDIUM
// ============================================================
// Find the minimum length contiguous subarray with sum >= target.
// Return 0 if no such subarray exists.
//
// Example: target=7, nums=[2,3,1,2,4,3] → 2  ([4,3])

// MinSubArrayLen returns the minimum length subarray with sum >= target.
// Time: O(n)  Space: O(1)
func MinSubArrayLen(target int, nums []int) int {
	// TODO: implement
	return 0
}
