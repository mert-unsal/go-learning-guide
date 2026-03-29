package stacks_queues

// PROBLEM 1: Valid Parentheses (LeetCode #20) — EASY
// Given a string of brackets '(', ')', '{', '}', '[', ']', return true if valid.
// Examples: "()" → true, "()[]{}" → true, "(]" → false, "([)]" → false
// Approach: push opening brackets onto stack. On closing, check top matches.
// Target: O(n) time, O(n) space.

func IsValid(s string) bool {
	// TODO: implement
	var stack []rune
	var enclosingValueMap = make(map[rune]rune)
	enclosingValueMap['('] = ')'
	enclosingValueMap[')'] = '('
	enclosingValueMap['{'] = '}'
	enclosingValueMap['}'] = '{'
	enclosingValueMap['['] = ']'
	enclosingValueMap[']'] = '['
	for _, ch := range s {
		if len(stack) == 0 {
			stack = append(stack, ch)
		} else if stack[len(stack)-1] == enclosingValueMap[ch] {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, ch)
		}
	}

	return len(stack) == 0
}
