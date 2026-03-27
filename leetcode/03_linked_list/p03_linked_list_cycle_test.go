package linked_list

import "testing"

func TestHasCycle(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		head := newList([]int{1, 2, 3})
		if HasCycle(head) {
			t.Error("HasCycle on acyclic list returned true")
		}
	})
	t.Run("has cycle", func(t *testing.T) {
		head := &ListNode{Val: 3}
		n2 := &ListNode{Val: 2}
		n3 := &ListNode{Val: 0}
		n4 := &ListNode{Val: -4}
		head.Next = n2
		n2.Next = n3
		n3.Next = n4
		n4.Next = n2 // cycle!
		if !HasCycle(head) {
			t.Error("HasCycle on cyclic list returned false")
		}
	})
	t.Run("single no cycle", func(t *testing.T) {
		head := &ListNode{Val: 1}
		if HasCycle(head) {
			t.Error("expected false")
		}
	})
}
