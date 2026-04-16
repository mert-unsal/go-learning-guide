package trees

// ============================================================
// PROBLEM 13: Same Tree (LeetCode #100) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the roots of two binary trees p and q, determine if they
//   are structurally identical and all corresponding nodes have the
//   same values.
//
// PARAMETERS:
//   p *TreeNode — root of the first binary tree (may be nil)
//   q *TreeNode — root of the second binary tree (may be nil)
//
// RETURN:
//   bool — true if both trees are structurally identical with equal values
//
// CONSTRAINTS:
//   • 0 <= number of nodes in each tree <= 100
//   • -10^4 <= Node.Val <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  p = [1,2,3], q = [1,2,3]
//   Output: true
//
// Example 2:
//   Input:  p = [1,2], q = [1,null,2]
//   Output: false
//   Why:    Structure differs — p has left child 2, q has right child 2
//
// Example 3:
//   Input:  p = [1,2,1], q = [1,1,2]
//   Output: false
//   Why:    Values at corresponding positions differ
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Recursive: both nil → true; one nil → false; vals differ → false; recurse both children
// • Iterative: BFS/DFS with paired node comparisons
// • Target: O(n) time, O(h) space where n = min(nodes in p, nodes in q)
func IsSameTree(p *TreeNode, q *TreeNode) bool {
	return false
}
