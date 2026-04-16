package trees

// ============================================================
// PROBLEM 11: Binary Tree Right Side View (LeetCode #199) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary tree, imagine yourself standing on
//   the right side of it. Return the values of the nodes you can
//   see ordered from top to bottom (the rightmost node at each level).
//
// PARAMETERS:
//   root *TreeNode — root of the binary tree (may be nil)
//
// RETURN:
//   []int — values visible from the right side, one per level
//
// CONSTRAINTS:
//   • 0 <= number of nodes <= 100
//   • -100 <= Node.Val <= 100
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [1,2,3,null,5,null,4]
//   Output: [1,3,4]
//   Why:    Level 0→1, Level 1→3 (rightmost), Level 2→4 (rightmost)
//
// Example 2:
//   Input:  root = [1,null,3]
//   Output: [1,3]
//
// Example 3:
//   Input:  root = []
//   Output: []
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • BFS: last node processed at each level is the right-side view node
// • DFS: visit right subtree first; add node if depth == len(result)
// • Target: O(n) time, O(h) space (DFS) or O(w) space (BFS, w = max width)
func RightSideView(root *TreeNode) []int {
	return nil
}
