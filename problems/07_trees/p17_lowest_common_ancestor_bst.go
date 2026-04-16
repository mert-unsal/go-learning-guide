package trees

// ============================================================
// PROBLEM 17: Lowest Common Ancestor of a BST (LeetCode #235) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given a binary search tree (BST) and two nodes p and q, find
//   their lowest common ancestor. The BST property allows you to
//   exploit value comparisons to navigate directly toward the LCA
//   without exploring the entire tree.
//
// PARAMETERS:
//   root *TreeNode — root of the BST
//   p    *TreeNode — first target node (guaranteed to exist)
//   q    *TreeNode — second target node (guaranteed to exist)
//
// RETURN:
//   *TreeNode — the lowest common ancestor node
//
// CONSTRAINTS:
//   • 2 <= number of nodes <= 10^5
//   • -10^9 <= Node.Val <= 10^9
//   • All Node.Val are unique
//   • p != q
//   • Both p and q exist in the BST
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [6,2,8,0,4,7,9,null,null,3,5], p = 2, q = 8
//   Output: 6
//   Why:    2 is in left subtree, 8 is in right subtree → root 6 is LCA
//
// Example 2:
//   Input:  root = [6,2,8,0,4,7,9,null,null,3,5], p = 2, q = 4
//   Output: 2
//   Why:    4 is in the subtree of 2, so 2 is the LCA (ancestor of itself)
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • If both p and q are less than root, LCA is in left subtree
// • If both p and q are greater than root, LCA is in right subtree
// • Otherwise, root is the split point → root is the LCA
// • Can be done iteratively (no recursion needed) for O(1) extra space
// • Target: O(h) time, O(1) space (iterative) or O(h) space (recursive)
func LowestCommonAncestorBST(root, p, q *TreeNode) *TreeNode {
	return nil
}
