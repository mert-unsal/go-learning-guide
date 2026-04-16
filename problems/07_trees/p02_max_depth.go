package trees

// ============================================================
// PROBLEM 2: Maximum Depth of Binary Tree (LeetCode #104) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, return its maximum depth.
//   The maximum depth is the number of nodes along the longest
//   path from the root node down to the farthest leaf node.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (may be nil)
//
// RETURN:
//   int — maximum depth of the tree (0 if root is nil)
//
// CONSTRAINTS:
//   • 0 <= number of nodes <= 10^4
//   • -100 <= Node.Val <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [3,9,20,null,null,15,7]
//   Output: 3
//   Why:    Longest path: 3→20→15 or 3→20→7, both length 3
//
// Example 2:
//   Input:  root = [1,null,2]
//   Output: 2
//
// Example 3:
//   Input:  root = []
//   Output: 0
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Recursive: depth = 1 + max(depth(left), depth(right)), nil returns 0
// • Iterative BFS: count levels using a queue
// • Target: O(n) time, O(h) space

func MaxDepth(root *TreeNode) int {
	return 0
}
