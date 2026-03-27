package dynamic_prog

// PROBLEM 1: Climbing Stairs (LeetCode #70) — EASY
// ways(n) = ways(n-1) + ways(n-2). Same as Fibonacci.
// Target: O(n) time, O(1) space.

func ClimbStairs(n int) int { return 0 }

// PROBLEM 2: Coin Change (LeetCode #322) — MEDIUM
// Fewest coins to make amount. dp[i] = min coins for amount i. Return -1 if impossible.
// Target: O(amount * len(coins)) time, O(amount) space.

func CoinChange(coins []int, amount int) int { return 0 }

// PROBLEM 3: House Robber (LeetCode #198) — MEDIUM
// Max amount robbing non-adjacent houses. rob(i) = max(rob(i-1), rob(i-2)+nums[i]).
// Target: O(n) time, O(1) space.

func Rob(nums []int) int { return 0 }

func max(a, b int) int { return 0 }

// PROBLEM 4: Unique Paths (LeetCode #62) — MEDIUM
// Robot on m×n grid, moves right/down. How many unique paths?
// dp[r][c] = dp[r-1][c] + dp[r][c-1].
// Target: O(m*n) time, O(n) space.

func UniquePaths(m int, n int) int { return 0 }

// PROBLEM 5: Longest Common Subsequence (LeetCode #1143) — MEDIUM
// LCS of two strings. dp[i][j] based on character match/mismatch.
// Target: O(m*n) time, O(m*n) space.

func LongestCommonSubsequence(text1 string, text2 string) int { return 0 }

// PROBLEM 6: Min Cost Climbing Stairs (LeetCode #746) — EASY
func MinCostClimbingStairs(cost []int) int { return 0 }
func min2(a, b int) int                    { return 0 }

// PROBLEM 7: Longest Increasing Subsequence (LeetCode #300) — MEDIUM
func LengthOfLIS(nums []int) int { return 0 }

// PROBLEM 8: Word Break (LeetCode #139) — MEDIUM
func WordBreak(s string, wordDict []string) bool { return false }

// PROBLEM 9: Jump Game (LeetCode #55) — MEDIUM
func CanJump(nums []int) bool { return false }

// PROBLEM 10: Partition Equal Subset Sum (LeetCode #416) — MEDIUM
func CanPartition(nums []int) bool { return false }

// PROBLEM 11: Maximum Product Subarray (LeetCode #152) — MEDIUM
func MaxProduct(nums []int) int { return 0 }

// PROBLEM 12: Decode Ways (LeetCode #91) — MEDIUM
func NumDecodings(s string) int { return 0 }

// PROBLEM 13: House Robber II (LeetCode #213) — MEDIUM
func RobII(nums []int) int                    { return 0 }
func robRange(nums []int, start, end int) int { return 0 }

// PROBLEM 14: Coin Change II (LeetCode #518) — MEDIUM
func CoinChangeII(amount int, coins []int) int { return 0 }

// PROBLEM 15: Interleaving String (LeetCode #97) — MEDIUM
func IsInterleave(s1, s2, s3 string) bool { return false }
