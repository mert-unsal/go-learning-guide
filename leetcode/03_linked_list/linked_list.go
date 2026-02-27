// Package linked_list contains LeetCode linked list problems with explanations.
// Topics: pointer manipulation, dummy head node, Floyd's cycle detection.
package linked_list

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

// ============================================================
// PROBLEM 1: Reverse Linked List (LeetCode #206) — EASY
// ============================================================
// Reverse a singly-linked list and return the new head.
//
// Example: 1→2→3→4→5 → 5→4→3→2→1
//
// Approach: three-pointer iterative.
// prev=nil, cur=head. Each step: save cur.Next, point cur.Next to prev,
// advance prev and cur. When cur=nil, prev is the new head.

// ReverseList reverses a linked list iteratively.
// Time: O(n)  Space: O(1)
func ReverseList(head *ListNode) *ListNode {
	var prev *ListNode
	cur := head

	for cur != nil {
		next := cur.Next // save next before we overwrite it
		cur.Next = prev  // reverse the pointer
		prev = cur       // advance prev
		cur = next       // advance cur
	}
	return prev // prev is now the new head
}

// ReverseListRecursive reverses a linked list recursively.
// Time: O(n)  Space: O(n) — recursion stack
// The base case: a list of 0 or 1 nodes is already reversed.
func ReverseListRecursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	// Recursively reverse the rest of the list
	newHead := ReverseListRecursive(head.Next)
	// Make the next node point back to head
	head.Next.Next = head
	head.Next = nil
	return newHead
}

// ============================================================
// PROBLEM 2: Merge Two Sorted Lists (LeetCode #21) — EASY
// ============================================================
// Merge two sorted linked lists and return the sorted merged list.
//
// Example: 1→2→4  and  1→3→4  →  1→1→2→3→4→4
//
// Approach: dummy head node eliminates edge cases for the first node.
// Compare the heads of both lists, attach the smaller one, advance that pointer.

// MergeTwoLists merges two sorted linked lists.
// Time: O(n + m)  Space: O(1) — we reuse existing nodes
func MergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	// Dummy head: we never have to special-case the first node
	dummy := &ListNode{}
	cur := dummy

	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			cur.Next = list1
			list1 = list1.Next
		} else {
			cur.Next = list2
			list2 = list2.Next
		}
		cur = cur.Next
	}

	// Attach the remaining non-nil list
	if list1 != nil {
		cur.Next = list1
	} else {
		cur.Next = list2
	}

	return dummy.Next // skip the dummy node
}

// ============================================================
// PROBLEM 3: Linked List Cycle (LeetCode #141) — EASY
// ============================================================
// Return true if the linked list has a cycle.
//
// Approach: Floyd's Tortoise and Hare algorithm.
// Slow pointer moves 1 step, fast pointer moves 2 steps.
// If there's a cycle, fast will eventually lap slow and they'll meet.
// If there's no cycle, fast will reach nil.

// HasCycle returns true if the list contains a cycle.
// Time: O(n)  Space: O(1)
func HasCycle(head *ListNode) bool {
	slow, fast := head, head

	for fast != nil && fast.Next != nil {
		slow = slow.Next      // move 1 step
		fast = fast.Next.Next // move 2 steps
		if slow == fast {
			return true // they met → cycle exists
		}
	}
	return false // fast reached end → no cycle
}

// ============================================================
// PROBLEM 4: Remove Nth Node From End (LeetCode #19) — MEDIUM
// ============================================================
// Remove the nth node from the end of the list.
//
// Example: 1→2→3→4→5, n=2 → 1→2→3→5
//
// Approach: two-pointer with n-step gap.
// Advance fast pointer n steps ahead. Then advance both until fast reaches end.
// Slow's next is the node to remove.

// RemoveNthFromEnd removes the nth node from the end.
// Time: O(n)  Space: O(1)
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	fast, slow := dummy, dummy

	// Advance fast n+1 steps (one extra so slow stops before the target)
	for i := 0; i <= n; i++ {
		fast = fast.Next
	}

	// Move both until fast reaches nil
	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}

	// slow.Next is the node to remove
	slow.Next = slow.Next.Next
	return dummy.Next
}

// ============================================================
// PROBLEM 5: Middle of the Linked List (LeetCode #876) — EASY
// ============================================================
// Return the middle node. If two middles exist, return the second.
//
// Example: 1→2→3→4→5 → node 3
// Example: 1→2→3→4   → node 3
//
// Approach: slow/fast pointers. When fast reaches end, slow is at middle.

// MiddleNode returns the middle node of the list.
// Time: O(n)  Space: O(1)
func MiddleNode(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}
