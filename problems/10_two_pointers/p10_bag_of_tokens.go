package two_pointers

// ============================================================
// PROBLEM 10: Bag of Tokens (LeetCode #948) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   You start with an initial power and a score of 0, and a bag of
//   tokens where each token has a value. You can play each token at
//   most once and in any order. With each token you can either:
//     - Play it face-up: if your power ≥ token value, lose that much
//       power and gain 1 score.
//     - Play it face-down: if your score ≥ 1, gain that much power
//       and lose 1 score.
//   Return the maximum score you can achieve.
//
// PARAMETERS:
//   tokens []int — array of token values
//   power  int   — initial power
//
// RETURN:
//   int — maximum achievable score
//
// CONSTRAINTS:
//   • 0 ≤ len(tokens) ≤ 1000
//   • 0 ≤ tokens[i], power ≤ 10⁴
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  tokens = [100], power = 50
//   Output: 0
//   Why:    Cannot play the only token face-up (50 < 100)
//
// Example 2:
//   Input:  tokens = [200, 100], power = 150
//   Output: 1
//   Why:    Play token 1 (100) face-up → power=50, score=1
//
// Example 3:
//   Input:  tokens = [100, 200, 300, 400], power = 200
//   Output: 2
//   Why:    Play 100 face-up (score=1), play 400 face-down (power=500),
//           play 200 face-up (score=2), play 300 face-up (score=3)? Actually max=2
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Sort tokens — buy cheapest (face-up) and sell most expensive (face-down)
// • Two pointers: left buys score with power, right sells score for power
// • Greedily buy from left; when stuck, sell from right if it gains net benefit
// • Track maxScore separately — don't return final score (it may dip after selling)
// • Target: O(n log n) time, O(1) extra space
func BagOfTokensScore(tokens []int, power int) int {
	return 0
}
