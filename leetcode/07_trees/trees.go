// Package trees contains LeetCode binary tree problems with explanations.
// Topics: DFS (recursive + iterative), BFS (level-order), LCA, tree construction.
package trees

// TreeNode is a binary tree node, matching LeetCode's definition.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// newTree builds a binary tree from level-order slice (0 = null node).
func newTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == 0 {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != 0 {
			node.Left = &TreeNode{Val: vals[i]}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != 0 {
			node.Right = &TreeNode{Val: vals[i]}
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}

// ============================================================
// PROBLEM 1: Binary Tree Inorder Traversal (LeetCode #94) — EASY
// ============================================================
// Return inorder traversal (Left → Root → Right).
//
// Example:     1
//               \
//                2
//               /
//              3
// Inorder: [1, 3, 2]

// InorderTraversal returns inorder traversal (Left → Root → Right).
// Time: O(n)  Space: O(n) for result + O(h) recursion stack
func InorderTraversal(root *TreeNode) []int {
	var result []int
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}
		dfs(node.Left)                    // visit left
		result = append(result, node.Val) // visit root
		dfs(node.Right)                   // visit right
	}
	dfs(root)
	return result
}

// InorderIterative uses a stack to avoid recursion.
func InorderIterative(root *TreeNode) []int {
	var result []int
	stack := []*TreeNode{}
	cur := root

	for cur != nil || len(stack) > 0 {
		// Push all left children
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		// Process node
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, cur.Val)
		cur = cur.Right // move to right subtree
	}
	return result
}

// ============================================================
// PROBLEM 2: Maximum Depth of Binary Tree (LeetCode #104) — EASY
// ============================================================
// Return the maximum depth (number of nodes along longest root-to-leaf path).
//
// Approach: DFS. depth = 1 + max(depth(left), depth(right))

// MaxDepth returns the maximum depth of the tree.
// Time: O(n)  Space: O(h) where h is tree height
func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	leftDepth := MaxDepth(root.Left)
	rightDepth := MaxDepth(root.Right)
	if leftDepth > rightDepth {
		return 1 + leftDepth
	}
	return 1 + rightDepth
}

// ============================================================
// PROBLEM 3: Level Order Traversal (LeetCode #102) — MEDIUM
// ============================================================
// Return values level by level (BFS). Each level is a separate slice.
//
// Approach: BFS with a queue. Process all nodes at current level before next.

// LevelOrder returns the level-order traversal grouped by level.
// Time: O(n)  Space: O(n)
func LevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	var result [][]int
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue)
		var level []int

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, level)
	}
	return result
}

// ============================================================
// PROBLEM 4: Lowest Common Ancestor (LeetCode #236) — MEDIUM
// ============================================================
// Find the lowest common ancestor of two nodes p and q in a binary tree.
// LCA is the deepest node that has both p and q as descendants.
//
// Key insight: recurse on both sides.
// If current node is p or q, return it.
// If both sides return non-nil, current node is LCA.
// Otherwise return whichever side found something.

// LowestCommonAncestor returns the LCA of nodes with values pVal and qVal.
// Time: O(n)  Space: O(h)
func LowestCommonAncestor(root *TreeNode, p *TreeNode, q *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root == p || root == q {
		return root // found one of the targets
	}

	left := LowestCommonAncestor(root.Left, p, q)
	right := LowestCommonAncestor(root.Right, p, q)

	if left != nil && right != nil {
		return root // p is in one subtree, q in the other → LCA is root
	}
	if left != nil {
		return left // both are in the left subtree
	}
	return right
}

// ============================================================
// PROBLEM 5: Symmetric Tree (LeetCode #101) — EASY
// ============================================================
// Return true if the tree is a mirror of itself.
//
// Approach: compare left and right subtrees recursively.
// A mirror: outer values match AND inner values match.

// IsSymmetric returns true if the tree is symmetric.
// Time: O(n)  Space: O(h)
func IsSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return isMirror(root.Left, root.Right)
}

func isMirror(left, right *TreeNode) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return left.Val == right.Val &&
		isMirror(left.Left, right.Right) && // outer pair
		isMirror(left.Right, right.Left) // inner pair
}

// ============================================================
// PROBLEM 6: Validate Binary Search Tree (LeetCode #98) — MEDIUM
// ============================================================
// Return true if the tree is a valid BST.
// Valid BST: left subtree < node < right subtree (ALL nodes, not just direct children)
//
// Approach: pass min/max bounds down. At each node, value must be within (min, max).

// IsValidBST returns true if the tree is a valid BST.
// Time: O(n)  Space: O(h)
func IsValidBST(root *TreeNode) bool {
	return isValid(root, nil, nil) // no bounds initially
}

func isValid(node *TreeNode, min, max *int) bool {
	if node == nil {
		return true
	}
	if min != nil && node.Val <= *min {
		return false // value must be GREATER than min (left bound)
	}
	if max != nil && node.Val >= *max {
		return false // value must be LESS than max (right bound)
	}
	v := node.Val
	return isValid(node.Left, min, &v) && // left: max bound tightens
		isValid(node.Right, &v, max) // right: min bound tightens
}

// ============================================================
// PROBLEM 7: Path Sum (LeetCode #112) — EASY
// ============================================================
// Return true if the tree has a root-to-leaf path whose values sum to targetSum.
//
// Example: root=[5,4,8,11,null,13,4,7,2,null,null,null,1], targetSum=22 → true
//
// Approach: DFS, subtract current node value from targetSum.
// At a leaf, check if remaining == 0.

// HasPathSum returns true if any root-to-leaf path sums to targetSum.
// Time: O(n)  Space: O(h)
func HasPathSum(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}
	targetSum -= root.Val
	if root.Left == nil && root.Right == nil {
		return targetSum == 0 // leaf: check remaining sum
	}
	return HasPathSum(root.Left, targetSum) || HasPathSum(root.Right, targetSum)
}

// ============================================================
// PROBLEM 8: Invert Binary Tree (LeetCode #226) — EASY
// ============================================================
// Invert (mirror) a binary tree.
//
// Example:
//     4          4
//    / \   →    / \
//   2   7      7   2
//  / \ / \    / \ / \
// 1  3 6  9  9  6 3  1

// InvertTree returns the root of the inverted tree.
// Time: O(n)  Space: O(h)
func InvertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	root.Left, root.Right = InvertTree(root.Right), InvertTree(root.Left)
	return root
}

// ============================================================
// PROBLEM 9: Diameter of Binary Tree (LeetCode #543) — EASY
// ============================================================
// The diameter is the length of the longest path between ANY two nodes
// (the path may not pass through the root).
//
// Approach: DFS. At each node, the path through it = leftHeight + rightHeight.
// Track the max across all nodes.

// DiameterOfBinaryTree returns the diameter (longest path edge count).
// Time: O(n)  Space: O(h)
func DiameterOfBinaryTree(root *TreeNode) int {
	maxDiameter := 0
	var height func(node *TreeNode) int
	height = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		l := height(node.Left)
		r := height(node.Right)
		if l+r > maxDiameter {
			maxDiameter = l + r
		}
		if l > r {
			return 1 + l
		}
		return 1 + r
	}
	height(root)
	return maxDiameter
}

// ============================================================
// PROBLEM 10: Construct Binary Tree from Preorder and Inorder (LeetCode #105) — MEDIUM
// ============================================================
// Reconstruct a binary tree from its preorder and inorder traversals.
//
// Example: preorder=[3,9,20,15,7], inorder=[9,3,15,20,7]
//          → tree with root 3
//
// Key insight: preorder[0] is always the root.
// Find that root in inorder to split left and right subtrees.

// BuildTree constructs a binary tree from preorder and inorder traversals.
// Time: O(n)  Space: O(n) for the index map
func BuildTree(preorder []int, inorder []int) *TreeNode {
	// Build index map for inorder lookup
	inorderIdx := make(map[int]int)
	for i, v := range inorder {
		inorderIdx[v] = i
	}

	var build func(preL, preR, inL, inR int) *TreeNode
	build = func(preL, preR, inL, inR int) *TreeNode {
		if preL > preR {
			return nil
		}
		rootVal := preorder[preL]
		root := &TreeNode{Val: rootVal}
		mid := inorderIdx[rootVal] // index in inorder
		leftSize := mid - inL
		root.Left = build(preL+1, preL+leftSize, inL, mid-1)
		root.Left = build(preL+1, preL+leftSize, inL, mid-1)
		root.Right = build(preL+leftSize+1, preR, mid+1, inR)
		return root
	}
	return build(0, len(preorder)-1, 0, len(inorder)-1)
}

// ============================================================
// PROBLEM 11: Binary Tree Right Side View (LeetCode #199) — MEDIUM
// ============================================================
// Return the values visible from the right side of the tree (one per level).
//
// Example: root=[1,2,3,null,5,null,4] → [1,3,4]
//
// Approach: BFS level-order; take the last element of each level.

// RightSideView returns values visible from the right of each level.
// Time: O(n)  Space: O(n)
func RightSideView(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var result []int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			if i == levelSize-1 {
				result = append(result, node.Val) // rightmost of this level
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
	return result
}

// ============================================================
// PROBLEM 12: Kth Smallest Element in a BST (LeetCode #230) — MEDIUM
// ============================================================
// Find the kth smallest value in a BST.
//
// Key insight: inorder traversal of a BST gives values in sorted order.
// Stop when we've visited k nodes.

// KthSmallest returns the kth smallest value in a BST.
// Time: O(H + k) where H is tree height  Space: O(H)
func KthSmallest(root *TreeNode, k int) int {
	count := 0
	result := 0
	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil || count >= k {
			return
		}
		inorder(node.Left)
		count++
		if count == k {
			result = node.Val
			return
		}
		inorder(node.Right)
	}
	inorder(root)
	return result
}
