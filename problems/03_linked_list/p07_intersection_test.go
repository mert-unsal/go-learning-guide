package linked_list

import "testing"

func TestGetIntersectionNode(t *testing.T) {
	// Create two lists that intersect
	shared := &ListNode{Val: 8, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5}}}
	headA := &ListNode{Val: 4, Next: &ListNode{Val: 1, Next: shared}}
	headB := &ListNode{Val: 5, Next: &ListNode{Val: 6, Next: &ListNode{Val: 1, Next: shared}}}

	got := GetIntersectionNode(headA, headB)
	if got != shared {
		t.Errorf("GetIntersectionNode: got %v, want node with val 8", got)
	}

	// No intersection
	a := newList([]int{1, 2, 3})
	b := newList([]int{4, 5, 6})
	got2 := GetIntersectionNode(a, b)
	if got2 != nil {
		t.Errorf("GetIntersectionNode: got %v, want nil", got2)
	}
}
