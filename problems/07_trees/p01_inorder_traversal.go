package trees

// ============================================================
// PROBLEM 1: Binary Tree Inorder Traversal (LeetCode #94) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, return the inorder traversal
//   of its nodes' values. Inorder: Left → Root → Right.
//   For a BST, this yields values in sorted (ascending) order.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (may be nil)
//
// RETURN:
//   []int — node values in inorder sequence
//
// CONSTRAINTS:
//   • 0 <= number of nodes <= 100
//   • -100 <= Node.Val <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [1,null,2,3]
//   Output: [1,3,2]
//   Why:    Left of 1 is nil → visit 1; right subtree has 3 left of 2 → [1,3,2]
//
// Example 2:
//   Input:  root = []
//   Output: []
//
// Example 3:
//   Input:  root = [1]
//   Output: [1]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Recursive: traverse left, append root, traverse right
// • Iterative: use an explicit stack — push left chain, pop and visit, go right
// • Morris traversal achieves O(1) space by temporarily threading the tree
// • Target: O(n) time, O(h) space (recursive/iterative), O(1) space (Morris)

func InorderTraversal(root *TreeNode) []int {
	return nil
}

// InorderIterative uses a stack instead of recursion.
func InorderIterative(root *TreeNode) []int {
	return nil
}
