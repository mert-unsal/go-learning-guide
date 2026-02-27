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
