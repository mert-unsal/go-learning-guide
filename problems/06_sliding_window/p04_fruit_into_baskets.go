package sliding_window

// ============================================================
// PROBLEM 4: Fruit Into Baskets (LeetCode #904) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   You are visiting a farm with a row of fruit trees. Each tree
//   produces one type of fruit (given as an integer). You have two
//   baskets, and each basket can only hold one type of fruit. Starting
//   from any tree, pick fruits moving right — you must stop when you
//   encounter a third distinct fruit type. Return the maximum number
//   of fruits you can collect.
//
// PARAMETERS:
//   fruits []int — an array where fruits[i] is the type of fruit at tree i
//
// RETURN:
//   int — the maximum number of fruits you can collect with two baskets
//
// CONSTRAINTS:
//   • 1 <= len(fruits) <= 10^5
//   • 0 <= fruits[i] < len(fruits)
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  fruits = [1,2,1]
//   Output: 3
//   Why:    All trees have at most 2 types — pick all 3 fruits.
//
// Example 2:
//   Input:  fruits = [0,1,2,2]
//   Output: 3
//   Why:    Pick trees [1,2,2] → types {1,2}, total 3 fruits.
//
// Example 3:
//   Input:  fruits = [1,2,3,2,2]
//   Output: 4
//   Why:    Pick trees [2,3,2,2] → types {2,3}, total 4 fruits.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • This is "longest subarray with at most 2 distinct values".
// • Use a variable-size sliding window with a frequency map. Shrink
//   the left side whenever distinct types exceed 2.
// • Target: O(n) time, O(1) space (map has at most 3 entries)

func TotalFruit(fruits []int) int {
	return 0
}
