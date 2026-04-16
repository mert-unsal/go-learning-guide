package trees

// ============================================================
// PROBLEM 7: Path Sum (LeetCode #112) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree and an integer targetSum,
//   return true if the tree has a root-to-leaf path such that
//   the sum of all node values along the path equals targetSum.
//   A leaf is a node with no children.
//
// PARAMETERS:
//   root      *TreeNode — root of the binary tree (may be nil)
//   targetSum int       — the target sum to find
//
// RETURN:
//   bool — true if a root-to-leaf path with the given sum exists
//
// CONSTRAINTS:
//   • 0 <= number of nodes <= 5000
//   • -1000 <= Node.Val <= 1000
//   • -1000 <= targetSum <= 1000
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [5,4,8,11,null,13,4,7,2,null,null,null,1], targetSum = 22
//   Output: true
//   Why:    Path 5→4→11→2 sums to 22
//
// Example 2:
//   Input:  root = [1,2,3], targetSum = 5
//   Output: false
//
// Example 3:
//   Input:  root = [], targetSum = 0
//   Output: false
//   Why:    Empty tree has no root-to-leaf path
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Subtract node value from targetSum as you recurse down
// • At a leaf, check if remaining targetSum equals the leaf's value
// • Be careful: must reach a leaf node — internal nodes don't count
// • Target: O(n) time, O(h) space
func HasPathSum(root *TreeNode, targetSum int) bool {
	return false
}
