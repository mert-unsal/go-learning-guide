package linked_list

// ============================================================
// Shared Types and Helpers for Linked List Problems
// ============================================================

// ListNode is a singly-linked list node, matching LeetCode's definition.
type ListNode struct {
	Val  int
	Next *ListNode
}

// newList is a test helper: builds a list from a slice.
func newList(vals []int) *ListNode {
	return nil
}

// toSlice is a test helper: converts a list to a slice.
func toSlice(head *ListNode) []int {
	return nil
}

// RandomNode is a node with an extra random pointer (for LeetCode #138).
type RandomNode struct {
	Val    int
	Next   *RandomNode
	Random *RandomNode
}
