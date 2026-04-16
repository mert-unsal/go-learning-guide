package trees

// ============================================================
// PROBLEM 18: Count Good Nodes in Binary Tree (LeetCode #1448) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, a node X is "good" if in the
//   path from the root to X there are no nodes with a value greater
//   than X. Return the number of good nodes in the binary tree.
//   The root is always a good node.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (at least one node)
//
// RETURN:
//   int — count of good nodes in the tree
//
// CONSTRAINTS:
//   • 1 <= number of nodes <= 10^5
//   • -10^4 <= Node.Val <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [3,1,4,3,null,1,5]
//   Output: 4
//   Why:    Good nodes: 3 (root), 3 (path max=3, 3>=3), 4 (path max=3, 4>=3), 5 (path max=4, 5>=4)
//
// Example 2:
//   Input:  root = [3,3,null,4,2]
//   Output: 3
//   Why:    Good nodes: 3 (root), 3 (3>=3), 4 (4>=3). Node 2 is not good (2<3)
//
// Example 3:
//   Input:  root = [1]
//   Output: 1
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • DFS passing the maximum value seen so far along the path from root
// • If current node's value >= maxSoFar, it's a good node — increment count
// • Update maxSoFar = max(maxSoFar, node.Val) when recursing into children
// • Target: O(n) time, O(h) space
func GoodNodes(root *TreeNode) int {
	return 0
}
