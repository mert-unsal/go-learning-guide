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
	return nil
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
