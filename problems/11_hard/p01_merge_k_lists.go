package hard

// ============================================================
// PROBLEM 1: Merge k Sorted Lists (LeetCode #23) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   You are given an array of k linked-lists, each sorted in
//   ascending order. Merge all the linked-lists into one sorted
//   linked-list and return it.
//
// PARAMETERS:
//   lists []*ListNode — slice of pointers to the heads of k sorted linked-lists
//
// RETURN:
//   *ListNode — head of the single merged sorted linked-list
//
// CONSTRAINTS:
//   • k == len(lists)
//   • 0 <= k <= 10^4
//   • 0 <= lists[i].length <= 500
//   • -10^4 <= lists[i][j] <= 10^4
//   • Each lists[i] is sorted in ascending order
//   • The sum of all lists[i].length will not exceed 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  lists = [[1,4,5],[1,3,4],[2,6]]
//   Output: [1,1,2,3,4,4,5,6]
//   Why:    All lists merged and sorted into one list.
//
// Example 2:
//   Input:  lists = []
//   Output: []
//
// Example 3:
//   Input:  lists = [[]]
//   Output: []
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Divide and conquer: recursively merge pairs of lists
// • Alternatively, use a min-heap to always pick the smallest head
// • Target: O(N log k) time, O(1) space (merge in-place) or O(k) with heap
func MergeKLists(lists []*ListNode) *ListNode {
	return nil
}
func mergeRange(lists []*ListNode, left, right int) *ListNode {
	return nil
}
func mergeTwoLists(l1, l2 *ListNode) *ListNode {
	return nil
}
