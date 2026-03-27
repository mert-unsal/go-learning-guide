package arrays

// InsertInterval ============================================================
// PROBLEM 16: Insert Interval (LeetCode #57) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//
//	You are given an array of NON-OVERLAPPING intervals sorted by
//	start time, and a new interval to insert. Insert the new interval
//	and merge if necessary. Return the result as a sorted list of
//	non-overlapping intervals.
//
//	Unlike Merge Intervals (#56), the input is ALREADY sorted and
//	non-overlapping. You just need to fit in one more interval.
//
// PARAMETERS:
//
//	intervals   [][]int — sorted non-overlapping intervals [start, end].
//	newInterval []int   — the interval to insert [start, end].
//
// RETURN:
//
//	[][]int — the result after insertion and merging.
//
// CONSTRAINTS:
//   - 0 <= intervals.length <= 10⁴
//   - intervals[i].length == 2
//   - 0 <= start_i <= end_i <= 10⁵
//   - intervals is sorted by start_i.
//   - newInterval.length == 2
//   - 0 <= newInterval[0] <= newInterval[1] <= 10⁵
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — New interval overlaps one:
//
//	Input:  intervals = [[1,3], [6,9]], newInterval = [2,5]
//	Output: [[1,5], [6,9]]
//	Why:    [2,5] overlaps with [1,3] → merge to [1,5]. [6,9] is untouched.
//
// Example 2 — New interval overlaps multiple:
//
//	Input:  intervals = [[1,2], [3,5], [6,7], [8,10], [12,16]], newInterval = [4,8]
//	Output: [[1,2], [3,10], [12,16]]
//	Why:    [4,8] overlaps with [3,5], [6,7], [8,10] → merge all to [3,10].
//
// Example 3 — New interval before all:
//
//	Input:  intervals = [[3,5], [6,9]], newInterval = [1,2]
//	Output: [[1,2], [3,5], [6,9]]
//
// Example 4 — New interval after all:
//
//	Input:  intervals = [[1,2], [3,4]], newInterval = [5,6]
//	Output: [[1,2], [3,4], [5,6]]
//
// Example 5 — Empty intervals:
//
//	Input:  intervals = [], newInterval = [5,7]
//	Output: [[5,7]]
//
// Example 6 — New interval contains all:
//
//	Input:  intervals = [[2,3], [4,5], [6,7]], newInterval = [1,10]
//	Output: [[1,10]]
//
// InsertInterval inserts a new interval and merges overlaps.
// Time: O(n)  Space: O(n)
func InsertInterval(intervals [][]int, newInterval []int) [][]int {
	// TODO: implement
	return nil
}
