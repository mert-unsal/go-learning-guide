package stacks_queues

// PROBLEM 5: Implement Queue using Stacks (LeetCode #232) — EASY
// Implement FIFO queue using only two stacks.
// Key insight: "inbox" and "outbox". Push to inbox. Pop/peek from outbox;
// if outbox empty, pour inbox into outbox (reverses order → FIFO).

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
