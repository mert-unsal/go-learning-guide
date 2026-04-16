package trees

// ============================================================
// PROBLEM 5: Symmetric Tree (LeetCode #101) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, check whether it is a mirror
//   of itself (i.e., symmetric around its center).
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (may be nil)
//
// RETURN:
//   bool — true if the tree is symmetric, false otherwise
//
// CONSTRAINTS:
//   • 1 <= number of nodes <= 1000
//   • -100 <= Node.Val <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [1,2,2,3,4,4,3]
//   Output: true
//   Why:    Left subtree [2,3,4] mirrors right subtree [2,4,3]
//
// Example 2:
//   Input:  root = [1,2,2,null,3,null,3]
//   Output: false
//   Why:    Left subtree has 3 on right, right subtree has 3 on right — not mirrored
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Use a helper isMirror(left, right) comparing two subtrees
// • Two trees mirror if roots equal, and left.Left mirrors right.Right, left.Right mirrors right.Left
// • Iterative: use a queue/stack with node pairs
// • Target: O(n) time, O(h) space
func IsSymmetric(root *TreeNode) bool {
	return false
}
func isMirror(left, right *TreeNode) bool {
	return false
}
