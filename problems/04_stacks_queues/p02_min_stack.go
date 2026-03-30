package stacks_queues

// PROBLEM 2: Min Stack (LeetCode #155) — MEDIUM
// Design a stack supporting push, pop, top, and getMin all in O(1).
// Key insight: auxiliary min-stack tracks the minimum at each level.

type MinStack struct {
	stack    []int
	minStack []int
}

func (s *MinStack) Push(val int) {
}
func (s *MinStack) Pop() {
}
func (s *MinStack) Top() int {
	return 0
}
func (s *MinStack) GetMin() int {
	return 0
}
