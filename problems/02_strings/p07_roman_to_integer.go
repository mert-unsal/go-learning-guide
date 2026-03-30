package strings_problems

// ============================================================
// PROBLEM 7: Roman to Integer (LeetCode #13) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Convert a Roman numeral string to an integer.
//   Roman numerals: I=1, V=5, X=10, L=50, C=100, D=500, M=1000.
//   Subtractive notation: IV=4, IX=9, XL=40, XC=90, CD=400, CM=900.
//
// CONSTRAINTS:
//   • 1 <= s.length <= 15
//   • s contains only characters: I, V, X, L, C, D, M.
//   • 1 <= answer <= 3999
//
// ─── EXAMPLES ───────────────────────────────────────────────
// Example 1: "III"     → 3
// Example 2: "LVIII"   → 58   (L=50, V=5, III=3)
// Example 3: "MCMXCIV" → 1994 (M=1000, CM=900, XC=90, IV=4)
// Example 4: "IV"      → 4    (subtractive: 5-1)
// Example 5: "IX"      → 9
//
// ─── THINGS TO THINK ABOUT ─────────────────────────────────
//   • If a smaller value appears BEFORE a larger one, subtract it.
//   • Otherwise, add it.
//   • Use a map from rune/byte → integer value.
//   • Target: O(n) time, O(1) space.

// RomanToInt converts a Roman numeral string to an integer.
// Time: O(n)  Space: O(1)
func RomanToInt(s string) int {
	return 0
}
