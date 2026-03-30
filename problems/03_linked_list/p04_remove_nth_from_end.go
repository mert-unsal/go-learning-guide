package linked_list

// ============================================================
// PROBLEM 4: Remove Nth Node From End (LeetCode #19) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the head of a linked list, remove the nth node from the
//   END of the list and return its head.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: [1,2,3,4,5], n=2 → [1,2,3,5]
// Example 2: [1], n=1 → []
// Example 3: [1,2], n=1 → [1]
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Two-pointer gap: advance fast n+1 ahead, then move both.
//   • When fast is nil, slow.Next is the node to remove.
//   • Use a dummy node to handle removing the head gracefully.
//   • Target: O(n) time, O(1) space.

// RemoveNthFromEnd removes the nth node from the end.
// Time: O(n)  Space: O(1)
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	return nil
}
