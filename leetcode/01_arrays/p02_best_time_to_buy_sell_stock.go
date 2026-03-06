package arrays

// ============================================================
// PROBLEM 2: Best Time to Buy and Sell Stock (LeetCode #121) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given an array prices where prices[i] is the price of a
//   given stock on the ith day.
//
//   You want to maximize your profit by choosing a SINGLE day to buy
//   one stock and choosing a DIFFERENT day IN THE FUTURE to sell it.
//
//   Return the maximum profit you can achieve from this transaction.
//   If you cannot achieve any profit, return 0.
//
//   Key rule: you MUST buy before you sell — you cannot sell on day 2
//   and buy on day 5.
//
// PARAMETERS:
//   prices []int — daily stock prices, one entry per day.
//
// RETURN:
//   int — the maximum profit achievable, or 0 if no profit is possible.
//
// CONSTRAINTS:
//   • 1 <= prices.length <= 10⁵
//   • 0 <= prices[i] <= 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Normal profit:
//   Input:  prices = [7, 1, 5, 3, 6, 4]
//   Output: 5
//   Why:    Buy on day 2 (price=1), sell on day 5 (price=6). Profit = 6 - 1 = 5.
//           Note: buying on day 2 and selling on day 1 is NOT allowed.
//
// Example 2 — Prices always decreasing (no profit):
//   Input:  prices = [7, 6, 4, 3, 1]
//   Output: 0
//   Why:    Every future price is lower than every past price.
//           No transaction yields a positive profit.
//
// Example 3 — Single day:
//   Input:  prices = [5]
//   Output: 0
//   Why:    You can't buy AND sell on the same day — need at least two days.
//
// Example 4 — Two days, profitable:
//   Input:  prices = [1, 5]
//   Output: 4
//   Why:    Buy day 1, sell day 2. Profit = 5 - 1 = 4.
//
// Example 5 — Minimum is not the best buy point:
//   Input:  prices = [2, 4, 1]
//   Output: 2
//   Why:    Buy day 1 (price=2), sell day 2 (price=4). Profit = 2.
//           The minimum price is on day 3, but there's no future day to sell.
//
// Example 6 — Multiple peaks and valleys:
//   Input:  prices = [3, 1, 4, 8, 7, 2, 5]
//   Output: 7
//   Why:    Buy on day 2 (price=1), sell on day 4 (price=8). Profit = 7.
//
// Example 7 — All same price:
//   Input:  prices = [3, 3, 3, 3]
//   Output: 0
//   Why:    No price difference means no profit.
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • What if you track the MINIMUM price seen so far as you scan?
//   • At each day, what is the profit if you sold today?
//   • You only need one pass through the array.
//   • Target: O(n) time, O(1) space.

// MaxProfit returns the maximum profit achievable from one transaction.
// Time: O(n)  Space: O(1)
func MaxProfit(prices []int) int {
	// TODO: implement
	return 0
}
