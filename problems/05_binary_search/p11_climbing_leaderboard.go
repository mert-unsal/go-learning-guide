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
	// Deduplicate ranked scores (already sorted descending)
	unique := []int{ranked[0]}
	for i := 1; i < len(ranked); i++ {
		if ranked[i] != ranked[i-1] {
			unique = append(unique, ranked[i])
		}
	}
	n := len(unique)
	result := make([]int, len(player))

	for i, score := range player {
		// Binary search in DESCENDING array: find first index where unique[mid] <= score
		// Rank = position in 1-based index + 1 for scores strictly above
		lo, hi := 0, n-1
		rank := n + 1 // default: after everyone
		for lo <= hi {
			mid := lo + (hi-lo)/2
			if unique[mid] <= score {
				rank = mid + 1 // player ties or beats unique[mid]
				hi = mid - 1   // try to find a better (earlier) position
			} else {
				lo = mid + 1
			}
		}
		result[i] = rank
	}
	return result
}
