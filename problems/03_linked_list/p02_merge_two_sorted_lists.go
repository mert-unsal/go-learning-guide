package linked_list

// ============================================================
// PROBLEM 2: Merge Two Sorted Lists (LeetCode #21) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Merge two sorted linked lists and return the head of the merged
//   sorted list. The list should be made by splicing together the
//   nodes of the two input lists.
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: [1,2,4] + [1,3,4] → [1,1,2,3,4,4]
// Example 2: [] + [] → []
// Example 3: [] + [0] → [0]
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • Use a DUMMY HEAD node to avoid edge cases for the first node.
//   • Compare heads, attach the smaller, advance that pointer.
//   • When one list is exhausted, attach the remaining other list.
//   • Target: O(n+m) time, O(1) space (reuse existing nodes).

// MergeTwoLists merges two sorted linked lists.
// Time: O(n + m)  Space: O(1)
func MergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	return nil
}
