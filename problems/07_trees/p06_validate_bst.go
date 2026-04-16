package trees

// ============================================================
// PROBLEM 6: Validate Binary Search Tree (LeetCode #98) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, determine if it is a valid
//   binary search tree (BST). A valid BST requires that for every
//   node, all values in its left subtree are strictly less than the
//   node's value, and all values in its right subtree are strictly
//   greater.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree
//
// RETURN:
//   bool — true if the tree is a valid BST, false otherwise
//
// CONSTRAINTS:
//   • 1 <= number of nodes <= 10^4
//   • -2^31 <= Node.Val <= 2^31 - 1
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [2,1,3]
//   Output: true
//   Why:    1 < 2 < 3, BST property holds for all nodes
//
// Example 2:
//   Input:  root = [5,1,4,null,null,3,6]
//   Output: false
//   Why:    Node 4 is in right subtree of 5 but 4 < 5 violates BST
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Pass min/max bounds down the tree; each node must be within (min, max)
// • Use *int (nil pointers) for initial unbounded min/max, or math.MinInt/MaxInt
// • Alternative: inorder traversal must produce strictly increasing sequence
// • Target: O(n) time, O(h) space
func IsValidBST(root *TreeNode) bool {
	return false
}
func isValid(node *TreeNode, min, max *int) bool {
	return false
}
