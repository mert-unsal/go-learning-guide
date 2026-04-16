package stacks_queues

// ============================================================
// PROBLEM 3: Daily Temperatures (LeetCode #739) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given an array of integers temperatures representing daily
//   temperatures, return an array answer such that answer[i] is the
//   number of days you have to wait after the i-th day to get a
//   warmer temperature. If there is no future day with a warmer
//   temperature, set answer[i] = 0.
//
// PARAMETERS:
//   temperatures []int — daily temperature readings
//
// RETURN:
//   []int — days until a warmer temperature for each day (0 if none)
//
// CONSTRAINTS:
//   • 1 <= len(temperatures) <= 10^5
//   • 30 <= temperatures[i] <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  temperatures = [73,74,75,71,69,72,76,73]
//   Output: [1,1,4,2,1,1,0,0]
//   Why:    Day 0 (73) waits 1 day for 74. Day 3 (71) waits 2 days for 72.
//
// Example 2:
//   Input:  temperatures = [30,40,50,60]
//   Output: [1,1,1,0]
//   Why:    Each day is warmer than the last; the final day has no warmer day.
//
// Example 3:
//   Input:  temperatures = [30,60,90]
//   Output: [1,1,0]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use a monotonic decreasing stack that stores indices.
// • When the current temp is warmer than the stack top, pop and
//   record the difference in indices as the answer.
// • Target: O(n) time, O(n) space

func DailyTemperatures(temperatures []int) []int {
	return nil
}
