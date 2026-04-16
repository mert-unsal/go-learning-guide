package hard

// ============================================================
// PROBLEM 12: Maximum Profit in Job Scheduling (LeetCode #1235) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   You have n jobs where every job is scheduled to be done from
//   startTime[i] to endTime[i], obtaining a profit of profit[i].
//   You can select a set of non-overlapping jobs to maximize total
//   profit. Return the maximum profit you can obtain.
//
// PARAMETERS:
//   startTime []int — start time of each job
//   endTime   []int — end time of each job
//   profit    []int — profit of each job
//
// RETURN:
//   int — maximum profit from a set of non-overlapping jobs
//
// CONSTRAINTS:
//   • 1 <= len(startTime) == len(endTime) == len(profit) <= 5 * 10^4
//   • 1 <= startTime[i] < endTime[i] <= 10^9
//   • 1 <= profit[i] <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  startTime = [1,2,3,3], endTime = [3,4,5,6], profit = [50,10,40,70]
//   Output: 120
//   Why:    Pick job 0 (profit 50) and job 3 (profit 70) → 120.
//
// Example 2:
//   Input:  startTime = [1,2,3,4,6], endTime = [3,5,10,6,9], profit = [20,20,100,70,60]
//   Output: 150
//   Why:    Pick job 0, job 3, and job 4 → 20 + 70 + 60 = 150.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort jobs by endTime, then DP: dp[i] = max profit using first i jobs
// • For each job, binary search for the latest non-overlapping job
// • dp[i] = max(dp[i-1], profit[i] + dp[lastNonOverlapping])
// • Target: O(n log n) time, O(n) space
func JobScheduling(startTime []int, endTime []int, profit []int) int {
	return 0
}
