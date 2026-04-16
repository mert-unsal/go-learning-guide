package dynamic_prog

// ============================================================
// PROBLEM 1b: Jumping on the Clouds (HackerRank) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   There is a series of clouds numbered 0..n-1. Each cloud is either
//   safe (0) or a thundercloud (1). You start on cloud 0 and must reach
//   cloud n-1. From any cloud i you can jump to i+1 or i+2, but you
//   must never land on a thundercloud. Return the minimum number of
//   jumps to reach the last cloud.
//
// PARAMETERS:
//   c []int — array of 0s and 1s representing clouds (c[0] and c[n-1] are always 0)
//
// RETURN:
//   int — minimum number of jumps to reach the last cloud
//
// CONSTRAINTS:
//   • 2 ≤ len(c) ≤ 100
//   • c[i] ∈ {0, 1}
//   • c[0] = 0 and c[n-1] = 0
//   • It is always possible to reach the last cloud
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  c = [0, 0, 1, 0, 0, 1, 0]
//   Output: 4
//   Why:    0→1→3→4→6 (skip thunderclouds at index 2 and 5)
//
// Example 2:
//   Input:  c = [0, 0, 0, 0, 1, 0]
//   Output: 3
//   Why:    0→2→3→5 (jump over thundercloud at index 4)
//
// Example 3:
//   Input:  c = [0, 0, 0, 1, 0, 0]
//   Output: 3
//   Why:    0→2→4→5
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Greedy: always try to jump 2 first; fall back to 1 if landing on thunder
// • No DP array needed — a single index pointer suffices
// • Target: O(n) time, O(1) space
func JumpingOnClouds(c []int) int {
	return 0
}
