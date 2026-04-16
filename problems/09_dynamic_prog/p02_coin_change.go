package dynamic_prog

// ============================================================
// PROBLEM 2: Coin Change (LeetCode #322) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of coin denominations and a total amount, return
//   the fewest number of coins needed to make up that amount. If the
//   amount cannot be made up by any combination of the coins, return -1.
//   You may assume an infinite supply of each coin denomination.
//
// PARAMETERS:
//   coins  []int — array of coin denominations (each ≥ 1)
//   amount int   — target amount (0 ≤ amount ≤ 10⁴)
//
// RETURN:
//   int — fewest coins to make amount, or -1 if impossible
//
// CONSTRAINTS:
//   • 1 ≤ len(coins) ≤ 12
//   • 1 ≤ coins[i] ≤ 2³¹ - 1
//   • 0 ≤ amount ≤ 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  coins = [1, 5, 10], amount = 12
//   Output: 3
//   Why:    10 + 1 + 1 = 12
//
// Example 2:
//   Input:  coins = [2], amount = 3
//   Output: -1
//   Why:    No combination of 2s can make 3
//
// Example 3:
//   Input:  coins = [1], amount = 0
//   Output: 0
//   Why:    Zero amount needs zero coins
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Bottom-up DP: dp[i] = min coins to make amount i
// • Recurrence: dp[i] = min(dp[i - coin] + 1) for each coin
// • Initialize dp with amount+1 (impossible sentinel) except dp[0] = 0
// • Target: O(amount × len(coins)) time, O(amount) space

func CoinChange(coins []int, amount int) int {
	return 0
}
