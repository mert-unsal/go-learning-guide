package heap_priority_queue

// ============================================================
// PROBLEM 4: Task Scheduler (LeetCode #621) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given an array of CPU tasks (each labeled with a
//   letter A-Z) and a cooling interval n. Each cycle the CPU
//   can complete one task or sit idle. Identical tasks must be
//   separated by at least n intervals. Return the minimum number
//   of intervals the CPU will take to finish all tasks.
//
// PARAMETERS:
//   tasks []byte — array of task labels ('A'-'Z')
//   n     int    — cooling interval between identical tasks
//
// RETURN:
//   int — minimum number of CPU intervals to complete all tasks
//
// CONSTRAINTS:
//   • 1 <= len(tasks) <= 10^4
//   • tasks[i] is an uppercase English letter
//   • 0 <= n <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  tasks = ['A','A','A','B','B','B'], n = 2
//   Output: 8
//   Why:    A → B → idle → A → B → idle → A → B (8 intervals).
//
// Example 2:
//   Input:  tasks = ['A','C','A','B','D','B'], n = 1
//   Output: 6
//   Why:    A → B → C → A → D → B — no idle needed.
//
// Example 3:
//   Input:  tasks = ['A','A','A','B','B','B'], n = 0
//   Output: 6
//   Why:    No cooldown, just run all 6 tasks.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Greedy with max-heap: always schedule the most frequent available task
// • Math formula: (maxFreq - 1) * (n + 1) + countOfMaxFreqTasks
// • Result = max(formula, len(tasks)) since all tasks must run
// • Target: O(n) time (counting sort), O(1) space (26 letters)
func LeastInterval(tasks []byte, n int) int {
	return 0
}
