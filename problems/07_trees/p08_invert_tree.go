package trees

// ============================================================
// PROBLEM 8: Invert Binary Tree (LeetCode #226) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, invert the tree (mirror it)
//   and return its root. Every left child is swapped with its
//   corresponding right child at every level.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (may be nil)
//
// RETURN:
//   *TreeNode — root of the inverted tree
//
// CONSTRAINTS:
//   • 0 <= number of nodes <= 100
//   • -100 <= Node.Val <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [4,2,7,1,3,6,9]
//   Output: [4,7,2,9,6,3,1]
//   Why:    Swap children at every node: 2↔7, 1↔3, 6↔9
//
// Example 2:
//   Input:  root = [2,1,3]
//   Output: [2,3,1]
//
// Example 3:
//   Input:  root = []
//   Output: []
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Recursive: swap left and right children, then recurse on both
// • Iterative: BFS/DFS with a queue/stack, swap children at each node
// • This is an in-place mutation — the original tree is modified
// • Target: O(n) time, O(h) space
func InvertTree(root *TreeNode) *TreeNode {
	return nil
}
