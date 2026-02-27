// Package dynamic_prog contains LeetCode dynamic programming problems.
// Topics: memoization, tabulation, 1D/2D DP, classic DP patterns.
package dynamic_prog

// ============================================================
// PROBLEM 1: Climbing Stairs (LeetCode #70) — EASY
// ============================================================
// You are climbing a staircase with n steps.
// Each time you can climb 1 or 2 steps. In how many distinct ways can you reach the top?
//
// Example: n=3 → 3  (1+1+1, 1+2, 2+1)
//
// Recurrence: ways(n) = ways(n-1) + ways(n-2)   ← same as Fibonacci!
// Base cases: ways(1) = 1, ways(2) = 2
//
// Use O(1) space: only need the previous two values.

// ClimbStairs returns the number of distinct ways to climb n stairs.
// Time: O(n)  Space: O(1)
func ClimbStairs(n int) int {
	if n <= 2 {
		return n
	}
	prev2, prev1 := 1, 2
	for i := 3; i <= n; i++ {
		cur := prev1 + prev2
		prev2 = prev1
		prev1 = cur
	}
	return prev1
}

// ============================================================
// PROBLEM 2: Coin Change (LeetCode #322) — MEDIUM
// ============================================================
// Given coin denominations and an amount, find the fewest coins to make that amount.
// Return -1 if impossible.
//
// Example: coins=[1,2,5], amount=11 → 3 (5+5+1)
//
// This is an UNBOUNDED KNAPSACK problem.
// dp[i] = minimum coins to make amount i
// dp[0] = 0 (no coins needed for amount 0)
// dp[i] = min(dp[i - coin] + 1) for all coins where coin <= i
//
// Initialize dp with amount+1 (a value larger than any valid answer) as "infinity".

// CoinChange returns the minimum number of coins to make amount, or -1 if impossible.
// Time: O(amount * len(coins))  Space: O(amount)
func CoinChange(coins []int, amount int) int {
	inf := amount + 1 // larger than any valid answer; acts as "infinity"
	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = inf // initialize with "infinity"
	}
	dp[0] = 0 // base case: 0 coins needed for amount 0

	for i := 1; i <= amount; i++ {
		for _, coin := range coins {
			if coin <= i && dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
			}
		}
	}

	if dp[amount] == inf {
		return -1
	}
	return dp[amount]
}

// ============================================================
// PROBLEM 3: House Robber (LeetCode #198) — MEDIUM
// ============================================================
// You are a robber. Houses are in a line. Adjacent houses have alarms —
// you cannot rob two adjacent houses. Maximize the total amount robbed.
//
// Example: nums=[1,2,3,1] → 4 (rob house 0 + house 2 = 1+3)
//
// Recurrence: rob(i) = max(rob(i-1), rob(i-2) + nums[i])
//   Either skip house i (take rob(i-1))
//   Or rob house i (take rob(i-2) + nums[i])
//
// Use two rolling variables instead of a full array.

// Rob returns the maximum amount that can be robbed without robbing adjacent houses.
// Time: O(n)  Space: O(1)
func Rob(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}
	prev2 := nums[0]               // max loot up to house i-2
	prev1 := max(nums[0], nums[1]) // max loot up to house i-1

	for i := 2; i < len(nums); i++ {
		cur := max(prev1, prev2+nums[i])
		prev2 = prev1
		prev1 = cur
	}
	return prev1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ============================================================
// PROBLEM 4: Unique Paths (LeetCode #62) — MEDIUM
// ============================================================
// An m×n grid. Robot starts at top-left, wants to reach bottom-right.
// Can only move right or down. How many unique paths?
//
// dp[r][c] = paths to reach (r, c)
// dp[r][c] = dp[r-1][c] + dp[r][c-1]  (came from top or from left)
// dp[0][*] = 1, dp[*][0] = 1  (only one way to traverse the first row/column)
//
// Space optimization: use a single row, update left-to-right.

// UniquePaths returns the number of unique paths in an m×n grid.
// Time: O(m*n)  Space: O(n)
func UniquePaths(m int, n int) int {
	dp := make([]int, n)
	for i := range dp {
		dp[i] = 1 // first row: all 1s
	}
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			dp[c] += dp[c-1] // dp[c] was top, dp[c-1] is left
		}
	}
	return dp[n-1]
}

// ============================================================
// PROBLEM 5: Longest Common Subsequence (LeetCode #1143) — MEDIUM
// ============================================================
// Given two strings text1 and text2, return the length of their LCS.
// A subsequence doesn't need to be contiguous.
//
// Example: text1="abcde", text2="ace" → 3 ("ace")
//
// dp[i][j] = LCS of text1[0..i-1] and text2[0..j-1]
// If text1[i-1] == text2[j-1]: dp[i][j] = dp[i-1][j-1] + 1
// Else:                         dp[i][j] = max(dp[i-1][j], dp[i][j-1])

// LongestCommonSubsequence returns the LCS length of text1 and text2.
// Time: O(m*n)  Space: O(m*n) — can be optimized to O(n)
func LongestCommonSubsequence(text1 string, text2 string) int {
	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	return dp[m][n]
}

// ============================================================
// PROBLEM 6: Min Cost Climbing Stairs (LeetCode #746) — EASY
// ============================================================
// Each step has a cost. You can start from step 0 or 1, and can climb
// 1 or 2 steps at a time. Find the minimum cost to reach the top.
//
// Example: cost=[10,15,20] → 15
//
// dp[i] = min cost to reach step i
// dp[i] = cost[i] + min(dp[i-1], dp[i-2])

// MinCostClimbingStairs returns minimum cost to reach the top.
// Time: O(n)  Space: O(1)
func MinCostClimbingStairs(cost []int) int {
	n := len(cost)
	prev2, prev1 := cost[0], cost[1]
	for i := 2; i < n; i++ {
		curr := cost[i] + min2(prev1, prev2)
		prev2 = prev1
		prev1 = curr
	}
	return min2(prev1, prev2)
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ============================================================
// PROBLEM 7: Longest Increasing Subsequence (LeetCode #300) — MEDIUM
// ============================================================
// Return the length of the longest strictly increasing subsequence.
//
// Example: nums=[10,9,2,5,3,7,101,18] → 4  ([2,3,7,101])
//
// Approach: dp[i] = length of LIS ending at index i.
// dp[i] = max(dp[j] + 1) for all j < i where nums[j] < nums[i]
//
// O(n log n) approach: maintain a "tails" array using binary search.

// LengthOfLIS returns the length of the longest increasing subsequence.
// Time: O(n log n)  Space: O(n)
func LengthOfLIS(nums []int) int {
	tails := []int{} // tails[i] = smallest tail of all increasing subsequences of length i+1

	for _, num := range nums {
		// Binary search for first tail >= num
		left, right := 0, len(tails)
		for left < right {
			mid := left + (right-left)/2
			if tails[mid] < num {
				left = mid + 1
			} else {
				right = mid
			}
		}
		if left == len(tails) {
			tails = append(tails, num) // extend
		} else {
			tails[left] = num // replace
		}
	}
	return len(tails)
}

// ============================================================
// PROBLEM 8: Word Break (LeetCode #139) — MEDIUM
// ============================================================
// Given a string s and a word dictionary, return true if s can be
// segmented into dictionary words.
//
// Example: s="leetcode", wordDict=["leet","code"] → true
//
// dp[i] = true if s[0..i-1] can be segmented
// dp[i] = dp[j] && s[j..i-1] in wordDict for some j

// WordBreak returns true if s can be segmented into dictionary words.
// Time: O(n² * m) where m = avg word length  Space: O(n)
func WordBreak(s string, wordDict []string) bool {
	wordSet := make(map[string]bool)
	for _, w := range wordDict {
		wordSet[w] = true
	}
	n := len(s)
	dp := make([]bool, n+1)
	dp[0] = true // empty string is always segmentable

	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if dp[j] && wordSet[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}
	return dp[n]
}

// ============================================================
// PROBLEM 9: Jump Game (LeetCode #55) — MEDIUM
// ============================================================
// You are at index 0. Each element is the max jump length at that position.
// Return true if you can reach the last index.
//
// Example: nums=[2,3,1,1,4] → true
// Example: nums=[3,2,1,0,4] → false
//
// Greedy: track the furthest index reachable. If i > furthest, stuck.

// CanJump returns true if you can reach the last index.
// Time: O(n)  Space: O(1)
func CanJump(nums []int) bool {
	furthest := 0
	for i, jump := range nums {
		if i > furthest {
			return false // can't reach this index
		}
		if i+jump > furthest {
			furthest = i + jump
		}
	}
	return true
}

// ============================================================
// PROBLEM 10: Partition Equal Subset Sum (LeetCode #416) — MEDIUM
// ============================================================
// Determine if the array can be partitioned into two subsets with equal sum.
//
// Example: nums=[1,5,11,5] → true ([1,5,5] and [11])
//
// This is a 0/1 Knapsack problem.
// Target = totalSum / 2. Can we pick elements summing to target?
// dp[j] = true if sum j is achievable using some subset.

// CanPartition returns true if the array can be split into two equal-sum subsets.
// Time: O(n * sum)  Space: O(sum)
func CanPartition(nums []int) bool {
	total := 0
	for _, n := range nums {
		total += n
	}
	if total%2 != 0 {
		return false // odd total can't be split evenly
	}
	target := total / 2
	dp := make([]bool, target+1)
	dp[0] = true

	for _, num := range nums {
		// Traverse backwards to avoid using same element twice
		for j := target; j >= num; j-- {
			dp[j] = dp[j] || dp[j-num]
		}
	}
	return dp[target]
}
