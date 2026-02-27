// Package stacks_queues contains LeetCode stack and queue problems.
// Topics: stack-based validation, auxiliary stack, monotonic stack.
package stacks_queues

// ============================================================
// PROBLEM 1: Valid Parentheses (LeetCode #20) — EASY
// ============================================================
// Given a string of brackets '(', ')', '{', '}', '[', ']',
// return true if it is valid (brackets close in the correct order).
//
// Examples: "()" → true, "()[]{}" → true, "(]" → false, "([)]" → false
//
// Approach: stack-based matching.
// Push each opening bracket. On closing bracket, check if top of stack matches.
// Valid if stack is empty at the end.

// IsValid returns true if the bracket string is valid.
// Time: O(n)  Space: O(n)
func IsValid(s string) bool {
	stack := []rune{}

	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, ch := range s {
		if ch == '(' || ch == '[' || ch == '{' {
			stack = append(stack, ch) // push opening bracket
		} else {
			// Closing bracket: check if top of stack matches
			if len(stack) == 0 || stack[len(stack)-1] != pairs[ch] {
				return false
			}
			stack = stack[:len(stack)-1] // pop
		}
	}
	return len(stack) == 0 // valid only if all brackets are matched
}

// ============================================================
// PROBLEM 2: Min Stack (LeetCode #155) — MEDIUM
// ============================================================
// Design a stack that supports push, pop, top, and getMin in O(1).
//
// Key insight: use an auxiliary min-stack that tracks the minimum
// at each level. When we push x, also push min(x, currentMin).

// MinStack supports O(1) push, pop, top, and getMin.
type MinStack struct {
	stack    []int
	minStack []int // parallel stack tracking current minimum
}

// Push adds x to the stack.
func (s *MinStack) Push(val int) {
	s.stack = append(s.stack, val)
	if len(s.minStack) == 0 || val < s.minStack[len(s.minStack)-1] {
		s.minStack = append(s.minStack, val)
	} else {
		s.minStack = append(s.minStack, s.minStack[len(s.minStack)-1])
	}
}

// Pop removes the top element.
func (s *MinStack) Pop() {
	s.stack = s.stack[:len(s.stack)-1]
	s.minStack = s.minStack[:len(s.minStack)-1]
}

// Top returns the top element without removing it.
func (s *MinStack) Top() int {
	return s.stack[len(s.stack)-1]
}

// GetMin returns the minimum element in the stack in O(1).
func (s *MinStack) GetMin() int {
	return s.minStack[len(s.minStack)-1]
}

// ============================================================
// PROBLEM 3: Daily Temperatures (LeetCode #739) — MEDIUM
// ============================================================
// Given an array of temperatures, return an array where answer[i] is
// the number of days until a warmer temperature.
// If no future warmer day exists, answer[i] = 0.
//
// Example: temperatures=[73,74,75,71,69,72,76,73]
//          answer=       [ 1, 1, 4, 2, 1, 1, 0, 0]
//
// Approach: monotonic decreasing stack (stores indices, not values).
// When we find a warmer day, pop all stack entries that are cooler.
// The answer for each popped index = current index - popped index.

// DailyTemperatures returns the wait days until a warmer temperature.
// Time: O(n)  Space: O(n)
func DailyTemperatures(temperatures []int) []int {
	n := len(temperatures)
	answer := make([]int, n) // zero-initialized = default "never" answer
	stack := []int{}         // stores indices of temperatures (decreasing order)

	for i, temp := range temperatures {
		// Pop all indices whose temperature is less than current
		for len(stack) > 0 && temperatures[stack[len(stack)-1]] < temp {
			prevIdx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			answer[prevIdx] = i - prevIdx // days waited
		}
		stack = append(stack, i)
	}
	return answer
}

// ============================================================
// PROBLEM 4: Evaluate Reverse Polish Notation (LeetCode #150) — MEDIUM
// ============================================================
// Evaluate an expression in Reverse Polish Notation (postfix).
// Valid operators: +, -, *, /. Division truncates toward zero.
//
// Example: ["2","1","+","3","*"] → 9   (meaning (2+1)*3)
// Example: ["4","13","5","/","+"] → 6  (meaning 4+(13/5))

// EvalRPN evaluates a Reverse Polish Notation expression.
// Time: O(n)  Space: O(n)
func EvalRPN(tokens []string) int {
	stack := []int{}

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			// Pop two operands (b is on top, a is below)
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var result int
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				result = a / b // Go integer division truncates toward zero
			}
			stack = append(stack, result)
		default:
			// It's a number — parse and push
			num := 0
			neg := false
			start := 0
			if token[0] == '-' {
				neg = true
				start = 1
			}
			for _, ch := range token[start:] {
				num = num*10 + int(ch-'0')
			}
			if neg {
				num = -num
			}
			stack = append(stack, num)
		}
	}
	return stack[0]
}
