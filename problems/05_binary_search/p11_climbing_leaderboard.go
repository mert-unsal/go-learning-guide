package binary_search

// ============================================================
// Climbing the Leaderboard — [M]
// ============================================================
// Given a leaderboard of scores (descending, with dense ranking) and
// a list of a player's scores, find their rank after each game.
// Dense rank = no gaps (1st, 2nd, 2nd, 3rd, NOT 1st, 2nd, 2nd, 4th).
//
// Example: ranked=[100,100,50,40,40,20,10], player=[5,25,50,120]
// Ranks: [6, 4, 2, 1]
//
// Approach: deduplicate ranked, binary search for each player score.

// ClimbingLeaderboard returns the player's rank after each score.
// Time: O((n + m) log n)  Space: O(n)
func ClimbingLeaderboard(ranked []int, player []int) []int {
	return nil
}
