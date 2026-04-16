package trees

// ============================================================
// PROBLEM 3: Binary Tree Level Order Traversal (LeetCode #102) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, return the level order
//   traversal of its nodes' values (i.e., from left to right,
//   level by level). Each level's values are in a separate slice.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (may be nil)
//
// RETURN:
//   [][]int — 2D slice where each inner slice is one level's values
//
// CONSTRAINTS:
//   • 0 <= number of nodes <= 2000
//   • -1000 <= Node.Val <= 1000
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [3,9,20,null,null,15,7]
//   Output: [[3],[9,20],[15,7]]
//   Why:    Level 0: [3], Level 1: [9,20], Level 2: [15,7]
//
// Example 2:
//   Input:  root = [1]
//   Output: [[1]]
//
// Example 3:
//   Input:  root = []
//   Output: []
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • BFS with a queue: process all nodes at current level before moving to next
// • Track level size with len(queue) at the start of each iteration
// • Can also be done with DFS + depth parameter to place values in correct level
// • Target: O(n) time, O(n) space

func LevelOrder(root *TreeNode) [][]int {
	return nil
}
