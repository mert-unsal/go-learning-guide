package stacks_queues

// ============================================================
// PROBLEM 2: Min Stack (LeetCode #155) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Design a stack that supports push, pop, top, and retrieving
//   the minimum element, all in O(1) time.
//
// PARAMETERS:
//   Push(val int) — pushes the element val onto the stack
//   Pop()         — removes the element on the top of the stack
//   Top() int     — gets the top element of the stack
//   GetMin() int  — retrieves the minimum element in the stack
//
// RETURN:
//   Each method returns as described above.
//
// CONSTRAINTS:
//   • -2^31 <= val <= 2^31 - 1
//   • Methods pop, top, and getMin will always be called on non-empty stacks
//   • At most 3 * 10^4 calls will be made to push, pop, top, and getMin
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  Push(-2), Push(0), Push(-3), GetMin(), Pop(), Top(), GetMin()
//   Output: -3, 0, -2
//   Why:    After pushing -2,0,-3 the min is -3. Pop -3, top is 0, min is -2.
//
// Example 2:
//   Input:  Push(1), Push(1), Top(), GetMin()
//   Output: 1, 1
//   Why:    Both elements are 1, so top and min are equal.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Maintain an auxiliary min-stack that tracks the current minimum
//   at each level of the main stack.
// • On Push, push onto minStack only if val <= current min.
// • Target: O(1) time per operation, O(n) space

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
