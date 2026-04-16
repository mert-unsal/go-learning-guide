package sliding_window

// ============================================================
// PROBLEM 6: Maximum Points from Cards (LeetCode #1423) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   There are several cards arranged in a row, each with a point
//   value. In one step, you can take one card from the beginning or
//   end of the row. You must take exactly k cards. Return the maximum
//   score you can obtain.
//
// PARAMETERS:
//   cardPoints []int — point values of cards in a row
//   k          int   — the number of cards to take
//
// RETURN:
//   int — the maximum total points from picking k cards
//
// CONSTRAINTS:
//   • 1 <= len(cardPoints) <= 10^5
//   • 1 <= cardPoints[i] <= 10^4
//   • 1 <= k <= len(cardPoints)
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  cardPoints = [1,2,3,4,5,6,1], k = 3
//   Output: 12
//   Why:    Take [1] from left and [6,1] from right → 1+6+1 is only 8.
//           Better: take [5,6,1] from right → 12.
//
// Example 2:
//   Input:  cardPoints = [2,2,2], k = 2
//   Output: 4
//   Why:    Any two cards sum to 4.
//
// Example 3:
//   Input:  cardPoints = [9,7,7,9,7,7,9], k = 7
//   Output: 55
//   Why:    Take all cards.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • The cards NOT taken form a contiguous subarray of size (n - k).
//   Maximize score = totalSum - minimum window sum of size (n - k).
// • Use a fixed-size sliding window of size (n - k) to find the
//   minimum subarray sum.
// • Target: O(n) time, O(1) space

func MaxScore(cardPoints []int, k int) int {
	return 0
}
