package trees

// ============================================================
// PROBLEM 10: Construct Binary Tree from Preorder and Inorder (LeetCode #105) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given two integer arrays preorder and inorder where preorder is
//   the preorder traversal and inorder is the inorder traversal of
//   the same binary tree, construct and return the binary tree.
//
// PARAMETERS:
//   preorder []int — preorder traversal of the tree
//   inorder  []int — inorder traversal of the tree
//
// RETURN:
//   *TreeNode — root of the constructed binary tree
//
// CONSTRAINTS:
//   • 1 <= preorder.length <= 3000
//   • inorder.length == preorder.length
//   • -3000 <= preorder[i], inorder[i] <= 3000
//   • All values are unique
//   • Each value of inorder also appears in preorder
//   • preorder and inorder represent the same tree
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
//   Output: [3,9,20,null,null,15,7]
//   Why:    preorder[0]=3 is root; in inorder, left of 3 is [9], right is [15,20,7]
//
// Example 2:
//   Input:  preorder = [-1], inorder = [-1]
//   Output: [-1]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • First element of preorder is always the root
// • Find root in inorder to split left/right subtrees
// • Use a hashmap for O(1) inorder index lookups instead of linear search
// • Recurse with slice boundaries or index ranges — avoid copying slices
// • Target: O(n) time, O(n) space
func BuildTree(preorder []int, inorder []int) *TreeNode {
	return nil
}
