package linked_list

// ============================================================
// PROBLEM 3: Linked List Cycle (LeetCode #141) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given head, determine if the linked list has a cycle in it.
//   A cycle exists if some node can be reached again by continuously
//   following next pointers.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: [3,2,0,-4] with tail connecting to node index 1 → true
// Example 2: [1,2] with tail connecting to node index 0 → true
// Example 3: [1] with no cycle → false
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Floyd's Tortoise and Hare: slow moves 1, fast moves 2.
//   • If cycle exists, fast will eventually meet slow.
//   • If no cycle, fast reaches nil.
//   • Target: O(n) time, O(1) space.

// HasCycle returns true if the list contains a cycle.
// Time: O(n)  Space: O(1)
func HasCycle(head *ListNode) bool {
	// TODO: implement
	return false
}
