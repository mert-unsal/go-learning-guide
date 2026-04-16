package stacks_queues

// ============================================================
// PROBLEM 4: Evaluate Reverse Polish Notation (LeetCode #150) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Evaluate the value of an arithmetic expression in Reverse Polish
//   Notation (postfix). Valid operators are +, -, *, /. Each operand
//   may be an integer or another expression. Division between two
//   integers truncates toward zero.
//
// PARAMETERS:
//   tokens []string — RPN tokens (operands as number strings, operators as "+","-","*","/")
//
// RETURN:
//   int — the result of evaluating the RPN expression
//
// CONSTRAINTS:
//   • 1 <= len(tokens) <= 10^4
//   • tokens[i] is an operator or an integer in the range [-200, 200]
//   • The answer and all intermediate calculations fit in a 32-bit integer
//   • The input always represents a valid RPN expression
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  tokens = ["2","1","+","3","*"]
//   Output: 9
//   Why:    ((2 + 1) * 3) = 9
//
// Example 2:
//   Input:  tokens = ["4","13","5","/","+"]
//   Output: 6
//   Why:    (4 + (13 / 5)) = 4 + 2 = 6 (division truncates toward zero)
//
// Example 3:
//   Input:  tokens = ["10","6","9","3","+","-11","*","/","*","17","+","5","+"]
//   Output: 22
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use a stack of integers. Push numbers; on an operator, pop two
//   operands, apply the operation, and push the result.
// • Be careful with operand order: second popped is the left operand.
// • Target: O(n) time, O(n) space

func EvalRPN(tokens []string) int {
	return 0
}
