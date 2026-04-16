package heap_priority_queue

// ============================================================
// PROBLEM 3: Meeting Rooms II (LeetCode #253) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of meeting time intervals where
//   intervals[i] = [start_i, end_i], return the minimum number
//   of conference rooms required to hold all meetings.
//
// PARAMETERS:
//   intervals [][]int — list of [start, end] meeting intervals
//
// RETURN:
//   int — minimum number of conference rooms needed
//
// CONSTRAINTS:
//   • 1 <= len(intervals) <= 10^4
//   • 0 <= start_i < end_i <= 10^6
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  intervals = [[0,30],[5,10],[15,20]]
//   Output: 2
//   Why:    [0,30] overlaps with [5,10] → need 2 rooms.
//
// Example 2:
//   Input:  intervals = [[7,10],[2,4]]
//   Output: 1
//   Why:    No overlap, one room suffices.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Min-heap of end times: sort by start, push end times, pop when earliest end <= current start
// • Alternatively, separate sorted start/end arrays with two-pointer sweep
// • Heap size at any point = rooms in use; track the maximum
// • Target: O(n log n) time, O(n) space
func MinMeetingRooms(intervals [][]int) int {
	return 0
}
