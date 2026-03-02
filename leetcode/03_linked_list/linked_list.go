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
	// TODO: implement
	return nil
}

// ReverseListRecursive reverses a linked list recursively.
// Time: O(n)  Space: O(n) — recursion stack
// The base case: a list of 0 or 1 nodes is already reversed.
func ReverseListRecursive(head *ListNode) *ListNode {
	// TODO: implement
	return nil
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
	// TODO: implement
	return nil
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
	// TODO: implement
	return false
}

// ============================================================
// PROBLEM 4: Remove Nth Node From End (LeetCode #19) — MEDIUM
// ============================================================
// Remove the nth node from the end of the list.
//
// Example: 1→2→3→4→5, n=2 → 1→2→3→5
//
// Approach: two-pointer with n-step gap.
// Advance fast pointer n+1 steps ahead. Then advance both until fast reaches end.
// Slow's next is the node to remove. Use a dummy node for edge cases.

// RemoveNthFromEnd removes the nth node from the end.
// Time: O(n)  Space: O(1)
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	// TODO: implement
	return head
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
	// TODO: implement
	return head
}

// ============================================================
// PROBLEM 6: Palindrome Linked List (LeetCode #234) — EASY
// ============================================================
// Return true if the linked list is a palindrome.
//
// Example: 1→2→2→1 → true
//
// Approach: find middle, reverse second half, compare with first half.

// IsPalindrome returns true if the linked list is a palindrome.
// Time: O(n)  Space: O(1)
func IsPalindrome(head *ListNode) bool {
	// TODO: implement
	return false
}

// ============================================================
// PROBLEM 7: Intersection of Two Linked Lists (LeetCode #160) — EASY
// ============================================================
// Find the node at which two linked lists intersect. Return nil if no intersection.
//
// Example: A: a1→a2→c1→c2→c3, B: b1→b2→b3→c1→c2→c3 → c1
//
// Key insight: two pointers. When pointer A reaches end, restart at head of B.
// When pointer B reaches end, restart at head of A.
// They will meet at the intersection (or both reach nil if no intersection).

// GetIntersectionNode returns the intersection node of two lists.
// Time: O(m+n)  Space: O(1)
func GetIntersectionNode(headA, headB *ListNode) *ListNode {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 8: Add Two Numbers (LeetCode #2) — MEDIUM
// ============================================================
// Two non-empty linked lists represent non-negative integers in reverse order.
// Add the two numbers and return the sum as a linked list.
//
// Example: (2→4→3) + (5→6→4) → 7→0→8  (342 + 465 = 807)
//
// Hint: use a carry variable, iterate while either list or carry is non-zero.

// AddTwoNumbers adds two numbers represented as reversed linked lists.
// Time: O(max(m,n))  Space: O(max(m,n))
func AddTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 9: Copy List with Random Pointer (LeetCode #138) — MEDIUM
// ============================================================
// A linked list where each node has Val, Next, and Random (any node or nil).
// Deep copy the list.
//
// Approach: hash map from original node → cloned node.

// RandomNode is a node with an extra random pointer.
type RandomNode struct {
	Val    int
	Next   *RandomNode
	Random *RandomNode
}

// CopyRandomList deep copies a linked list with random pointers.
// Time: O(n)  Space: O(n)
func CopyRandomList(head *RandomNode) *RandomNode {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 10: Reorder List (LeetCode #143) — MEDIUM
// ============================================================
// Given a list L0→L1→...→Ln, reorder to: L0→Ln→L1→Ln-1→L2→Ln-2→...
//
// Example: 1→2→3→4→5 → 1→5→2→4→3
//
// Approach: find middle, reverse second half, merge two halves.

// ReorderList reorders the list in-place.
// Time: O(n)  Space: O(1)
func ReorderList(head *ListNode) {
	// TODO: implement
}

// ============================================================
// PROBLEM 11: Linked List Cycle II (LeetCode #142) — MEDIUM
// ============================================================
// Given a linked list, return the node where the cycle begins.
// Return nil if there is no cycle.
//
// Approach: Floyd's cycle detection. Once slow and fast meet,
// move one pointer back to head and advance both by 1.
// They will meet at the cycle entry point.
//
// Proof: let d = distance to cycle start, c = cycle length.
// When they meet, slow traveled d + k, fast traveled d + k + nc.
// Since fast = 2*slow: d + k + nc = 2(d + k) → d = nc - k.
// So moving d steps from head and d steps from meeting point both reach cycle start.

// DetectCycle returns the node where the cycle begins, or nil.
// Time: O(n)  Space: O(1)
func DetectCycle(head *ListNode) *ListNode {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 12: Reverse Nodes in k-Group (LeetCode #25) — HARD
// ============================================================
// Reverse the nodes of a linked list k at a time.
// If the number of nodes is not a multiple of k, the remaining nodes
// at the end stay as they are.
//
// Example: [1,2,3,4,5], k=2 → [2,1,4,3,5]
// Example: [1,2,3,4,5], k=3 → [3,2,1,4,5]

// ReverseKGroup reverses nodes in groups of k.
// Time: O(n)  Space: O(1)
func ReverseKGroup(head *ListNode, k int) *ListNode {
	// TODO: implement
	return head
}

// ============================================================
// PROBLEM 13: LRU Cache (LeetCode #146) — MEDIUM
// ============================================================
// Design a data structure that follows LRU eviction policy.
// get(key) and put(key, value) must both run in O(1).
//
// Approach: doubly-linked list + hash map.
// - Hash map: key → node pointer (O(1) lookup)
// - Doubly-linked list: most recently used at head, least recently at tail
// - On get: move node to head
// - On put: add to head; if over capacity, evict tail

type lruNode struct {
	key, val   int
	prev, next *lruNode
}

// LRUCache is an LRU cache with O(1) get and put.
type LRUCache struct {
	capacity   int
	cache      map[int]*lruNode
	head, tail *lruNode // sentinel nodes
}

// NewLRUCache creates a new LRU cache with given capacity.
func NewLRUCache(capacity int) *LRUCache {
	// TODO: implement — create head/tail sentinels, link them
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*lruNode),
	}
}

// Get returns the value for key, or -1 if not found. Marks as recently used.
func (c *LRUCache) Get(key int) int {
	// TODO: implement
	return -1
}

// Put inserts or updates a key-value pair. Evicts LRU entry if over capacity.
func (c *LRUCache) Put(key, value int) {
	// TODO: implement
}

// ============================================================
// PROBLEM 14: Find the Duplicate Number (LeetCode #287) — MEDIUM
// ============================================================
// Given an array of n+1 integers where each is in [1, n], find the
// duplicate number. Must not modify the array. O(1) extra space.
//
// Example: nums=[1,3,4,2,2] → 2
//
// Approach: Floyd's cycle detection on index mapping.
// Treat nums as a linked list: index i → nums[i].
// Since there's a duplicate, there must be a cycle. The cycle entry is the duplicate.

// FindDuplicate finds the duplicate number using Floyd's algorithm.
// Time: O(n)  Space: O(1)
func FindDuplicate(nums []int) int {
	// TODO: implement
	return 0
}
