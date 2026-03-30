package arrays

// MergeIntervals ============================================================
// PROBLEM 15: Merge Intervals (LeetCode #56) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//
//	Given an array of intervals where intervals[i] = [start_i, end_i],
//	merge all OVERLAPPING intervals and return an array of the
//	non-overlapping intervals that cover all the intervals in the input.
//
//	Two intervals overlap if one starts before the other ends:
//	  [1,3] and [2,6] overlap → merge to [1,6].
//	  [1,3] and [4,6] do NOT overlap (gap between 3 and 4).
//	  [1,4] and [4,6] DO overlap (they share the point 4) → merge to [1,6].
//
// PARAMETERS:
//
//	intervals [][]int — each element is [start, end] with start <= end.
//
// RETURN:
//
//	[][]int — merged non-overlapping intervals, sorted by start.
//
// CONSTRAINTS:
//   - 1 <= intervals.length <= 10⁴
//   - intervals[i].length == 2
//   - 0 <= start_i <= end_i <= 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Some overlapping:
//
//	Input:  intervals = [[1,3], [2,6], [8,10], [15,18]]
//	Output: [[1,6], [8,10], [15,18]]
//	Why:    [1,3] and [2,6] overlap → merged to [1,6].
//	        [8,10] and [15,18] don't overlap with anything.
//
// Example 2 — All overlapping:
//
//	Input:  intervals = [[1,4], [4,5]]
//	Output: [[1,5]]
//	Why:    They share the point 4 → merge to [1,5].
//
// Example 3 — No overlaps:
//
//	Input:  intervals = [[1,2], [4,5], [7,8]]
//	Output: [[1,2], [4,5], [7,8]]
//
// Example 4 — Single interval:
//
//	Input:  intervals = [[1,10]]
//	Output: [[1,10]]
//
// Example 5 — Unsorted input:
//
//	Input:  intervals = [[3,4], [1,2], [5,6]]
//	Output: [[1,2], [3,4], [5,6]]
//	Why:    Intervals are not necessarily given in sorted order!
//
// Example 6 — One interval contains another:
//
//	Input:  intervals = [[1,10], [3,5]]
//	Output: [[1,10]]
//	Why:    [3,5] is completely inside [1,10].
//
// Example 7 — Chain of overlaps:
//
//	Input:  intervals = [[1,3], [2,4], [3,5], [4,6]]
//	Output: [[1,6]]
//	Why:    Each overlaps with the next, forming one big merged interval.
//
// MergeIntervals merges overlapping intervals.
// Time: O(n log n)  Space: O(n)
func MergeIntervals(intervals [][]int) [][]int {
	return nil
}
