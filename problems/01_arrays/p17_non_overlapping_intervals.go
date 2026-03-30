package arrays

// EraseOverlapIntervals ============================================================
// PROBLEM 17: Non-overlapping Intervals (LeetCode #435) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//
//	Given an array of intervals intervals where intervals[i] = [start_i, end_i],
//	return the MINIMUM NUMBER of intervals you need to REMOVE to make the
//	rest of the intervals non-overlapping.
//
//	Two intervals [a, b) and [b, c) are considered NON-overlapping
//	(touching endpoints are fine).
//
// PARAMETERS:
//
//	intervals [][]int — each element is [start, end].
//
// RETURN:
//
//	int — the minimum number of intervals to remove.
//
// CONSTRAINTS:
//   - 1 <= intervals.length <= 10⁵
//   - intervals[i].length == 2
//   - -5 × 10⁴ <= start_i < end_i <= 5 × 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Remove one:
//
//	Input:  intervals = [[1,2], [2,3], [3,4], [1,3]]
//	Output: 1
//	Why:    Remove [1,3]. The remaining [[1,2],[2,3],[3,4]] don't overlap.
//
// Example 2 — Remove one:
//
//	Input:  intervals = [[1,2], [1,2], [1,2]]
//	Output: 2
//	Why:    Keep one [1,2], remove the other two duplicates.
//
// Example 3 — No removal needed:
//
//	Input:  intervals = [[1,2], [2,3]]
//	Output: 0
//	Why:    [1,2] and [2,3] don't overlap (touching at 2 is OK).
//
// Example 4 — All overlapping:
//
//	Input:  intervals = [[1,10], [2,3], [4,5], [6,7]]
//	Output: 1
//	Why:    Remove [1,10]. The rest [2,3],[4,5],[6,7] don't overlap.
//	        Removing 1 big interval is better than removing 3 small ones.
//
// Example 5 — Nested intervals:
//
//	Input:  intervals = [[1,5], [2,3]]
//	Output: 1
//	Why:    Either remove [1,5] or [2,3]. Removing [1,5] is better —
//	        it leaves [2,3] which is shorter.
//
// Example 6 — Chain of overlaps:
//
//	Input:  intervals = [[1,3], [2,4], [3,5]]
//	Output: 1
//	Why:    Remove [2,4]. Then [1,3] and [3,5] are non-overlapping.
//
// EraseOverlapIntervals returns the minimum removals for non-overlapping intervals.
// Time: O(n log n)  Space: O(1)
func EraseOverlapIntervals(intervals [][]int) int {
	return 0
}
