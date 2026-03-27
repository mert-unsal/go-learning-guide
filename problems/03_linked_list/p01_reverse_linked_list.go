package linked_list

// ============================================================
// PROBLEM 1: Reverse Linked List (LeetCode #206) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the head of a singly linked list, reverse the list, and
//   return the reversed list.
//
// PARAMETERS:
//   head *ListNode — the head of the linked list (may be nil).
//
// RETURN:
//   *ListNode — the head of the reversed list.
//
// CONSTRAINTS:
//   • The number of nodes is in [0, 5000].
//   • -5000 <= Node.val <= 5000.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1: [1,2,3,4,5] → [5,4,3,2,1]
// Example 2: [1,2] → [2,1]
// Example 3: [] → []
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Iterative: use three pointers (prev, curr, next).
//   • Recursive: the recursive call reverses the rest; then fix pointers.
//   • Target: O(n) time. Iterative = O(1) space, Recursive = O(n) stack.

// ReverseList reverses a linked list iteratively.
// Time: O(n)  Space: O(1)
func ReverseList(head *ListNode) *ListNode {
	// TODO: implement
	return nil
}

// ReverseListRecursive reverses a linked list recursively.
// Time: O(n)  Space: O(n) — recursion stack
func ReverseListRecursive(head *ListNode) *ListNode {
	// TODO: implement
	return nil
}
