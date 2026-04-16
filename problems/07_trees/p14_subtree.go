package trees

// ============================================================
// PROBLEM 14: Subtree of Another Tree (LeetCode #572) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the roots of two binary trees root and subRoot, return
//   true if there is a subtree of root with the same structure and
//   node values as subRoot. A subtree consists of a node in root
//   and all of that node's descendants.
//
// PARAMETERS:
//   root    *TreeNode — root of the main binary tree
//   subRoot *TreeNode — root of the candidate subtree
//
// RETURN:
//   bool — true if subRoot is a subtree of root
//
// CONSTRAINTS:
//   • 1 <= number of nodes in root <= 2000
//   • 1 <= number of nodes in subRoot <= 1000
//   • -10^4 <= Node.Val <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [3,4,5,1,2], subRoot = [4,1,2]
//   Output: true
//   Why:    The left subtree of root (rooted at 4) matches subRoot exactly
//
// Example 2:
//   Input:  root = [3,4,5,1,2,null,null,null,null,0], subRoot = [4,1,2]
//   Output: false
//   Why:    Node 2 has a left child 0 in root but not in subRoot
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • For each node in root, check if the subtree rooted there equals subRoot (use IsSameTree)
// • Recursive: isSubtree(root, sub) = isSame(root, sub) || isSubtree(root.Left, sub) || isSubtree(root.Right, sub)
// • Optimization: serialize both trees and use string matching (KMP) for O(m+n) time
// • Target: O(m*n) time naïve, O(m+n) with serialization; O(h) space
func IsSubtree(root *TreeNode, subRoot *TreeNode) bool {
	return false
}
