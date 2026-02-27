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

// ============================================================
// PROBLEM 5: Implement Queue using Stacks (LeetCode #232) — EASY
// ============================================================
// Implement a FIFO queue using only two stacks.
//
// Key insight: use an "inbox" and "outbox" stack.
// Push always goes to inbox. Pop/peek from outbox; if outbox is empty,
// pour all inbox items into outbox (this reverses the order → FIFO).

// MyQueue implements a queue using two stacks.
type MyQueue struct {
	inbox  []int
	outbox []int
}

// Push adds element to the back of the queue.
func (q *MyQueue) Push(x int) {
	q.inbox = append(q.inbox, x)
}

func (q *MyQueue) pour() {
	if len(q.outbox) == 0 {
		for len(q.inbox) > 0 {
			top := q.inbox[len(q.inbox)-1]
			q.inbox = q.inbox[:len(q.inbox)-1]
			q.outbox = append(q.outbox, top)
		}
	}
}

// Pop removes the front element and returns it.
func (q *MyQueue) Pop() int {
	q.pour()
	val := q.outbox[len(q.outbox)-1]
	q.outbox = q.outbox[:len(q.outbox)-1]
	return val
}

// Peek returns the front element without removing it.
func (q *MyQueue) Peek() int {
	q.pour()
	return q.outbox[len(q.outbox)-1]
}

// Empty returns true if the queue is empty.
func (q *MyQueue) Empty() bool {
	return len(q.inbox) == 0 && len(q.outbox) == 0
}

// ============================================================
// PROBLEM 6: Next Greater Element I (LeetCode #496) — EASY
// ============================================================
// For each element in nums1, find the next greater element in nums2.
// Return -1 if no greater element exists.
//
// Example: nums1=[4,1,2], nums2=[1,3,4,2] → [-1,3,-1]
//
// Approach: precompute next greater for all elements in nums2 using
// a monotonic decreasing stack, store in a map.

// NextGreaterElement returns next greater element for each nums1 value.
// Time: O(n+m)  Space: O(n)
func NextGreaterElement(nums1 []int, nums2 []int) []int {
	nextGreater := make(map[int]int)
	stack := []int{}

	for _, num := range nums2 {
		// Pop all elements smaller than current — current is their next greater
		for len(stack) > 0 && stack[len(stack)-1] < num {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nextGreater[top] = num
		}
		stack = append(stack, num)
	}
	// Remaining in stack have no next greater element
	for _, num := range stack {
		nextGreater[num] = -1
	}

	result := make([]int, len(nums1))
	for i, num := range nums1 {
		result[i] = nextGreater[num]
	}
	return result
}

// ============================================================
// PROBLEM 7: Decode String (LeetCode #394) — MEDIUM
// ============================================================
// Given an encoded string like "3[a]2[bc]", return "aaabcbc".
// Nested patterns are possible: "2[abc]3[cd]ef" → "abcabccdcdcdef"
//
// Approach: stack-based. Push counts and partial strings onto stacks.
// When ']' is encountered, pop and repeat.

// DecodeString decodes the encoded string.
// Time: O(output length)  Space: O(output length)
func DecodeString(s string) string {
	countStack := []int{}
	strStack := []string{}
	current := ""
	k := 0

	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			k = k*10 + int(ch-'0')
		} else if ch == '[' {
			countStack = append(countStack, k)
			strStack = append(strStack, current)
			current = ""
			k = 0
		} else if ch == ']' {
			count := countStack[len(countStack)-1]
			countStack = countStack[:len(countStack)-1]
			prev := strStack[len(strStack)-1]
			strStack = strStack[:len(strStack)-1]
			repeated := ""
			for i := 0; i < count; i++ {
				repeated += current
			}
			current = prev + repeated
		} else {
			current += string(ch)
		}
	}
	return current
}

// ============================================================
// PROBLEM 8: Largest Rectangle in Histogram (LeetCode #84) — HARD
// ============================================================
// Given bar heights, find the area of the largest rectangle.
//
// Example: heights=[2,1,5,6,2,3] → 10
//
// Approach: monotonic increasing stack stores indices.
// When a shorter bar is found, pop and compute area with popped bar as height.
// Width = current index - stack top - 1 (or full width if stack is empty).

// LargestRectangleArea returns the largest rectangle area in a histogram.
// Time: O(n)  Space: O(n)
func LargestRectangleArea(heights []int) int {
	stack := []int{} // indices, monotonic increasing by height
	maxArea := 0
	n := len(heights)

	for i := 0; i <= n; i++ {
		h := 0
		if i < n {
			h = heights[i]
		}
		for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
			height := heights[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]
			width := i
			if len(stack) > 0 {
				width = i - stack[len(stack)-1] - 1
			}
			area := height * width
			if area > maxArea {
				maxArea = area
			}
		}
		stack = append(stack, i)
	}
	return maxArea
}
