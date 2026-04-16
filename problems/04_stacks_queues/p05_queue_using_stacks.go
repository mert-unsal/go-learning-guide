package stacks_queues

// ============================================================
// PROBLEM 5: Implement Queue using Stacks (LeetCode #232) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Implement a first-in-first-out (FIFO) queue using only two stacks.
//   The queue should support push (to back), pop (from front), peek
//   (front element), and empty.
//
// PARAMETERS:
//   Push(x int) — pushes element x to the back of the queue
//   Pop() int   — removes and returns the element from the front
//   Peek() int  — returns the front element without removing it
//   Empty() bool — returns true if the queue is empty
//
// RETURN:
//   Each method returns as described above.
//
// CONSTRAINTS:
//   • 1 <= x <= 9
//   • At most 100 calls will be made to push, pop, peek, and empty
//   • All calls to pop and peek are valid (queue is non-empty)
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  Push(1), Push(2), Peek(), Pop(), Empty()
//   Output: 1, 1, false
//   Why:    Peek returns front (1), Pop removes 1, queue still has 2.
//
// Example 2:
//   Input:  Push(1), Pop(), Empty()
//   Output: 1, true
//   Why:    After popping the only element, queue is empty.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use an "inbox" stack and an "outbox" stack. Push always goes to inbox.
// • On Pop/Peek: if outbox is empty, pour all of inbox into outbox
//   (this reverses the order, giving FIFO behavior).
// • Target: O(1) amortized time per operation, O(n) space

type MyQueue struct {
	inbox  []int
	outbox []int
}

func (q *MyQueue) Push(x int) {
}
func (q *MyQueue) Pop() int {
	return 0
}
func (q *MyQueue) Peek() int {
	return 0
}
func (q *MyQueue) Empty() bool {
	return false
}
