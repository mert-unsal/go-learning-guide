package backtracking

// ============================================================
// PROBLEM 6: Letter Combinations of a Phone Number (LeetCode #17) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a string containing digits from 2-9 inclusive, return
//   all possible letter combinations that the number could
//   represent. Return the answer in any order. The digit-to-letter
//   mapping follows a telephone keypad (2=abc, 3=def, ..., 9=wxyz).
//
// PARAMETERS:
//   digits string — string of digits ('2'-'9')
//
// RETURN:
//   []string — all possible letter combinations
//
// CONSTRAINTS:
//   • 0 <= len(digits) <= 4
//   • digits[i] is a digit in the range ['2', '9']
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  digits = "23"
//   Output: ["ad","ae","af","bd","be","bf","cd","ce","cf"]
//   Why:    2→{a,b,c} × 3→{d,e,f} = 9 combinations.
//
// Example 2:
//   Input:  digits = ""
//   Output: []
//
// Example 3:
//   Input:  digits = "2"
//   Output: ["a","b","c"]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Backtracking: build combinations one digit at a time
// • Use a map/slice for digit→letters mapping
// • Iterative BFS-style: expand the result list one digit at a time
// • Target: O(4^n) time where n=len(digits), O(n) space for recursion
func LetterCombinations(digits string) []string {
	return nil
}
