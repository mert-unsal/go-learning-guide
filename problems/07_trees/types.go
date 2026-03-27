package trees

// Shared types and helpers for tree problems.

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

// TrieNode represents a node in the trie (for LeetCode #208).
type TrieNode struct {
	children [26]*TrieNode
	isEnd    bool
}

// Trie is a prefix tree for lowercase English words.
type Trie struct {
	root *TrieNode
}
