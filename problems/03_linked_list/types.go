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
	if len(vals) == 0 {
		return nil
	}
	head := &ListNode{Val: vals[0]}
	cur := head
	for _, v := range vals[1:] {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return head
}

// toSlice is a test helper: converts a list to a slice.
func toSlice(head *ListNode) []int {
	var res []int
	for cur := head; cur != nil; cur = cur.Next {
		res = append(res, cur.Val)
	}
	return res
}

// RandomNode is a node with an extra random pointer (for LeetCode #138).
type RandomNode struct {
	Val    int
	Next   *RandomNode
	Random *RandomNode
}
