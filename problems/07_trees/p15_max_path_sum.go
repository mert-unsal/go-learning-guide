package trees

// ============================================================
// PROBLEM 15: Binary Tree Maximum Path Sum (LeetCode #124) — HARD
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, return the maximum path sum
//   of any non-empty path. A path is a sequence of nodes where
//   each pair of adjacent nodes has an edge connecting them. A node
//   can only appear in the path at most once. The path does not
//   need to pass through the root.
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (at least one node)
//
// RETURN:
//   int — the maximum path sum across all possible paths
//
// CONSTRAINTS:
//   • 1 <= number of nodes <= 3 * 10^4
//   • -1000 <= Node.Val <= 1000
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [1,2,3]
//   Output: 6
//   Why:    Path 2→1→3 has sum 6, the optimal path
//
// Example 2:
//   Input:  root = [-10,9,20,null,null,15,7]
//   Output: 42
//   Why:    Path 15→20→7 has sum 42 (skipping the -10 root entirely)
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • At each node, compute max gain = node.Val + max(0, leftGain) + max(0, rightGain)
// • Update global max with the "split" path through the node
// • Return to parent only the best single-branch gain: node.Val + max(0, max(leftGain, rightGain))
// • Negative gains should be clamped to 0 (don't extend a path that hurts the sum)
// • Target: O(n) time, O(h) space
func MaxPathSum(root *TreeNode) int {
	return 0
}
