package stacks_queues

// ============================================================
// PROBLEM 10: Car Fleet (LeetCode #853) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   There are n cars going to the same destination along a one-lane
//   road. Each car i starts at position[i] with speed speed[i]. A
//   car can never pass another car ahead of it, but it can catch up
//   and drive at the same speed (forming a fleet). Return the number
//   of car fleets that will arrive at the destination.
//
// PARAMETERS:
//   target   int   — the destination position
//   position []int — starting position of each car
//   speed    []int — speed of each car
//
// RETURN:
//   int — the number of car fleets arriving at the target
//
// CONSTRAINTS:
//   • n == len(position) == len(speed)
//   • 1 <= n <= 10^5
//   • 0 < target <= 10^6
//   • 0 <= position[i] < target
//   • 0 < speed[i] <= 10^6
//   • All positions are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  target = 12, position = [10,8,0,5,3], speed = [2,4,1,1,3]
//   Output: 3
//   Why:    Car at 10 arrives at t=1. Cars at 8,5 merge (fleet). Car at 0,3
//           never catch up → 3 fleets.
//
// Example 2:
//   Input:  target = 10, position = [3], speed = [3]
//   Output: 1
//   Why:    Single car forms a single fleet.
//
// Example 3:
//   Input:  target = 100, position = [0,2,4], speed = [4,2,1]
//   Output: 1
//   Why:    All cars eventually merge into one fleet.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort cars by position descending. Calculate arrival time for each.
// • Iterate: if a car's arrival time > current fleet time, it starts a new fleet.
// • Target: O(n log n) time, O(n) space

func CarFleet(target int, position []int, speed []int) int {
	return 0
}
