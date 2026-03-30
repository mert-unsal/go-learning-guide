package linked_list

import "testing"

func TestCopyRandomList(t *testing.T) {
	// Basic test: single node with no random pointer
	node := &RandomNode{Val: 1}
	got := CopyRandomList(node)
	if got != nil && got == node {
		t.Errorf("CopyRandomList should return a deep copy, not the same node")
	}

	// Nil input
	got2 := CopyRandomList(nil)
	if got2 != nil {
		t.Errorf("CopyRandomList(nil) = %v, want nil", got2)
	}
}
