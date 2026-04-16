package trees

// ============================================================
// PROBLEM 4: Lowest Common Ancestor of a Binary Tree (LeetCode #236) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a binary tree and two nodes p and q, find their lowest
//   common ancestor (LCA). The LCA is the deepest node that is
//   an ancestor of both p and q (a node can be an ancestor of itself).
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree
//   p    *TreeNode — first target node (guaranteed to exist in tree)
//   q    *TreeNode — second target node (guaranteed to exist in tree)
//
// RETURN:
//   *TreeNode — the lowest common ancestor node
//
// CONSTRAINTS:
//   • 2 <= number of nodes <= 10^5
//   • -10^9 <= Node.Val <= 10^9
//   • All Node.Val are unique
//   • p != q
//   • Both p and q exist in the tree
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 1
//   Output: 3
//   Why:    Node 3 is the deepest node that is ancestor of both 5 and 1
//
// Example 2:
//   Input:  root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 4
//   Output: 5
//   Why:    Node 5 is ancestor of itself and of 4
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Recursive post-order: if current node is p or q, return it
// • If left and right subtrees both return non-nil, current node is the LCA
// • If only one side returns non-nil, propagate that result up
// • Target: O(n) time, O(h) space
func LowestCommonAncestor(root *TreeNode, p *TreeNode, q *TreeNode) *TreeNode {
	return nil
}
