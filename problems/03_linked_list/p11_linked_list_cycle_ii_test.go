package linked_list

import "testing"

func TestDetectCycle(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		head := newList([]int{1, 2, 3})
		if DetectCycle(head) != nil {
			t.Error("expected nil for acyclic list")
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
		n4.Next = n2 // cycle at n2
		got := DetectCycle(head)
		if got != n2 {
			t.Error("expected cycle start at node 2")
		}
	})
}
