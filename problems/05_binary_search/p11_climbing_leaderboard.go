package binary_search

// ============================================================
// PROBLEM 11: Climbing the Leaderboard (HackerRank) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a leaderboard of scores in descending order using dense
//   ranking (no gaps: 1st, 2nd, 2nd, 3rd — not 1st, 2nd, 2nd, 4th),
//   and a list of a player's game scores in ascending order, return
//   the player's rank after each game.
//
// PARAMETERS:
//   ranked []int — leaderboard scores in descending order (may have duplicates)
//   player []int — the player's scores in ascending order
//
// RETURN:
//   []int — the player's rank after each game score
//
// CONSTRAINTS:
//   • 1 <= len(ranked) <= 2 * 10^5
//   • 1 <= len(player) <= 2 * 10^5
//   • 0 <= ranked[i], player[i] <= 10^9
//   • ranked is sorted in non-increasing order
//   • player is sorted in non-decreasing order
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  ranked = [100,100,50,40,40,20,10], player = [5,25,50,120]
//   Output: [6,4,2,1]
//   Why:    Dense ranks are [1,1,2,3,3,4,5]. Score 5→rank 6, 25→4, 50→2, 120→1.
//
// Example 2:
//   Input:  ranked = [100,90,90,80], player = [70,80,105]
//   Output: [4,3,1]
//   Why:    Dense ranks [1,2,2,3]. Score 70→4, 80→3 (ties with 80), 105→1.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Deduplicate the ranked array to get unique scores with dense ranks.
// • For each player score, binary search the deduplicated array to
//   find the insertion position, which gives the rank.
// • Target: O((n + m) log n) time, O(n) space

// ClimbingLeaderboard returns the player's rank after each score.
// Time: O((n + m) log n)  Space: O(n)
func ClimbingLeaderboard(ranked []int, player []int) []int {
	return nil
}
