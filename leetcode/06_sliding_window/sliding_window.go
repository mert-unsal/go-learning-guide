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
	// Build the first window
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += nums[i]
	}
	maxSum := windowSum

	// Slide window: add right element, remove left element
	for i := k; i < len(nums); i++ {
		windowSum += nums[i] - nums[i-k]
		if windowSum > maxSum {
			maxSum = windowSum
		}
	}
	return float64(maxSum) / float64(k)
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
	if len(s) == 0 || len(t) == 0 {
		return ""
	}

	// Count required characters from t
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	required := len(need) // number of unique chars we need to satisfy

	have := make(map[byte]int)
	formed := 0 // how many unique chars currently satisfy their frequency

	left := 0
	minLen := len(s) + 1
	minLeft := 0

	for right := 0; right < len(s); right++ {
		ch := s[right]
		have[ch]++

		// Check if this char now satisfies its required count
		if count, ok := need[ch]; ok && have[ch] == count {
			formed++
		}

		// Try to shrink window from the left while it's valid
		for formed == required {
			// Update minimum window
			if right-left+1 < minLen {
				minLen = right - left + 1
				minLeft = left
			}
			// Remove leftmost character
			leftCh := s[left]
			have[leftCh]--
			if count, ok := need[leftCh]; ok && have[leftCh] < count {
				formed-- // window no longer satisfies this char
			}
			left++
		}
	}

	if minLen == len(s)+1 {
		return ""
	}
	return s[minLeft : minLeft+minLen]
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
	if len(s1) > len(s2) {
		return false
	}
	var need, have [26]int

	// Count s1 characters
	for i := 0; i < len(s1); i++ {
		need[s1[i]-'a']++
	}

	// Initialize window with first len(s1) chars of s2
	for i := 0; i < len(s1); i++ {
		have[s2[i]-'a']++
	}
	if need == have {
		return true
	}

	// Slide the window
	for i := len(s1); i < len(s2); i++ {
		have[s2[i]-'a']++         // add new right char
		have[s2[i-len(s1)]-'a']-- // remove old left char
		if need == have {
			return true
		}
	}
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
	basket := make(map[int]int) // fruit type → count in window
	left := 0
	maxFruits := 0

	for right := 0; right < len(fruits); right++ {
		basket[fruits[right]]++

		// Shrink window if more than 2 distinct types
		for len(basket) > 2 {
			leftFruit := fruits[left]
			basket[leftFruit]--
			if basket[leftFruit] == 0 {
				delete(basket, leftFruit)
			}
			left++
		}

		if right-left+1 > maxFruits {
			maxFruits = right - left + 1
		}
	}
	return maxFruits
}
