package trees

// ============================================================
// PROBLEM 12: Kth Smallest Element in a BST (LeetCode #230) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   Given the root of a binary search tree and an integer k,
//   return the kth smallest value (1-indexed) among all node
//   values in the BST.
//
// PARAMETERS:
//   root *TreeNode — root of the BST
//   k    int       — the 1-based rank of the desired element
//
// RETURN:
//   int — the kth smallest value in the BST
//
// CONSTRAINTS:
//   • 1 <= k <= number of nodes <= 10^4
//   • 0 <= Node.Val <= 10^4
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  root = [3,1,4,null,2], k = 1
//   Output: 1
//   Why:    Inorder traversal is [1,2,3,4]; 1st smallest is 1
//
// Example 2:
//   Input:  root = [5,3,6,2,4,null,null,1], k = 3
//   Output: 3
//   Why:    Inorder traversal is [1,2,3,4,5,6]; 3rd smallest is 3
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Inorder traversal of BST yields sorted order — stop at kth element
// • Iterative inorder with a stack avoids traversing the entire tree
// • Decrement k at each visit; when k reaches 0, that's the answer
// • Target: O(H + k) time, O(H) space where H is tree height
func KthSmallest(root *TreeNode, k int) int {
	return 0
}
