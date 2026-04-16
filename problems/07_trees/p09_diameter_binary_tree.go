package trees

// ============================================================
// PROBLEM 9: Diameter of Binary Tree (LeetCode #543) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, return the length of the
//   diameter. The diameter is the length of the longest path
//   between any two nodes (measured in number of edges). This
//   path may or may not pass through the root.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (may be nil)
//
// RETURN:
//   int — the diameter (number of edges on the longest path)
//
// CONSTRAINTS:
//   • 1 <= number of nodes <= 10^4
//   • -100 <= Node.Val <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [1,2,3,4,5]
//   Output: 3
//   Why:    Longest path is 4→2→1→3 or 5→2→1→3, both have 3 edges
//
// Example 2:
//   Input:  root = [1,2]
//   Output: 1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • At each node, diameter through it = height(left) + height(right)
// • Track the global maximum diameter while computing heights recursively
// • Use a closure variable or pass a pointer to track the running max
// • Target: O(n) time, O(h) space
func DiameterOfBinaryTree(root *TreeNode) int {
	return 0
}
