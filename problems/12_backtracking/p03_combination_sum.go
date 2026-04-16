package backtracking

// ============================================================
// PROBLEM 3: Combination Sum (LeetCode #39) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of distinct integers candidates and a target
//   integer target, return a list of all unique combinations of
//   candidates where the chosen numbers sum to target. The same
//   number may be chosen an unlimited number of times. Two
//   combinations are unique if the frequency of at least one chosen
//   number is different.
//
// PARAMETERS:
//   candidates []int — array of distinct positive integers
//   target     int   — target sum
//
// RETURN:
//   [][]int — all unique combinations that sum to target
//
// CONSTRAINTS:
//   • 1 <= len(candidates) <= 30
//   • 2 <= candidates[i] <= 40
//   • All elements of candidates are distinct
//   • 1 <= target <= 40
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  candidates = [2,3,6,7], target = 7
//   Output: [[2,2,3],[7]]
//   Why:    2+2+3=7 and 7=7 are the only combinations.
//
// Example 2:
//   Input:  candidates = [2,3,5], target = 8
//   Output: [[2,2,2,2],[2,3,3],[3,5]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Backtracking: at each step, choose candidate[i..n] (allow reuse)
// • Avoid duplicates by only considering candidates from current index onward
// • Prune when remaining target < current candidate
// • Target: O(n^(T/M)) time where T=target, M=min(candidates)
func CombinationSum(candidates []int, target int) [][]int {
	return nil
}

// ============================================================
// PROBLEM 4: Combination Sum II (LeetCode #40) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a collection of candidate numbers (candidates) and a
//   target number (target), find all unique combinations in
//   candidates where the candidate numbers sum to target. Each
//   number in candidates may only be used ONCE in the combination.
//   The solution set must not contain duplicate combinations.
//
// PARAMETERS:
//   candidates []int — array of integers (may contain duplicates)
//   target     int   — target sum
//
// RETURN:
//   [][]int — all unique combinations that sum to target (no duplicate combos)
//
// CONSTRAINTS:
//   • 1 <= len(candidates) <= 100
//   • 1 <= candidates[i] <= 50
//   • 1 <= target <= 30
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  candidates = [10,1,2,7,6,1,5], target = 8
//   Output: [[1,1,6],[1,2,5],[1,7],[2,6]]
//   Why:    Each combination uses elements at most once; no duplicate combos.
//
// Example 2:
//   Input:  candidates = [2,5,2,1,2], target = 5
//   Output: [[1,2,2],[5]]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort candidates first to group duplicates together
// • Backtracking: skip duplicate candidates at the same recursion depth
// • Key dedup logic: if candidates[i] == candidates[i-1] and i > start, skip
// • Target: O(2^n) time, O(n) space (excluding output)
func CombinationSum2(candidates []int, target int) [][]int {
	return nil
}
