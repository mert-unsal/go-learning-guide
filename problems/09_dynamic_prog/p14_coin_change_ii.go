package dynamic_prog

// ============================================================
// PROBLEM 14: Coin Change II (LeetCode #518) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of coin denominations and a total amount, return
//   the number of combinations that make up that amount. You have
//   an infinite supply of each coin. If the amount cannot be made,
//   return 0.
//
// PARAMETERS:
//   amount int   — target amount
//   coins  []int — array of coin denominations
//
// RETURN:
//   int — number of combinations that sum to amount
//
// CONSTRAINTS:
//   • 1 ≤ len(coins) ≤ 300
//   • 1 ≤ coins[i] ≤ 5000
//   • 0 ≤ amount ≤ 5000
//   • All values of coins are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  amount = 5, coins = [1, 2, 5]
//   Output: 4
//   Why:    5=5, 5=2+2+1, 5=2+1+1+1, 5=1+1+1+1+1
//
// Example 2:
//   Input:  amount = 3, coins = [2]
//   Output: 0
//   Why:    No combination of 2s can make 3
//
// Example 3:
//   Input:  amount = 10, coins = [10]
//   Output: 1
//   Why:    Only one coin equals the amount
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Unbounded knapsack: iterate coins in outer loop, amounts in inner loop
// • dp[j] += dp[j - coin] (order of loops matters to avoid counting permutations)
// • dp[0] = 1 (one way to make amount 0: use no coins)
// • Compare with Coin Change I: this counts combinations, not min coins
// • Target: O(amount × len(coins)) time, O(amount) space
func CoinChangeII(amount int, coins []int) int {
	return 0
}
