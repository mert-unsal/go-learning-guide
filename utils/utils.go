// Package utils provides shared data structures and helper functions
// used across the LeetCode practice problems in this project.
//
// LeetCode defines its own ListNode and TreeNode types in problems;
// here we mirror those definitions so our solutions feel authentic.
package utils

import "fmt"

// ============================================================
// LINKED LIST — ListNode
// ============================================================
// LeetCode always gives you: type ListNode struct { Val int; Next *ListNode }
// We define the same thing here.

// ListNode represents a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// NewList builds a linked list from a slice and returns the head.
// Example: NewList([]int{1,2,3}) → 1 → 2 → 3 → nil
func NewList(vals []int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	head := &ListNode{Val: vals[0]}
	cur := head
	for _, v := range vals[1:] {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return head
}

// ListToSlice converts a linked list back to a slice (useful in tests).
// Stops after 1000 nodes to avoid infinite loops in cyclic lists.
func ListToSlice(head *ListNode) []int {
	var result []int
	cur := head
	for cur != nil && len(result) < 1000 {
		result = append(result, cur.Val)
		cur = cur.Next
	}
	return result
}

// PrintList prints a linked list in a readable format: 1 -> 2 -> 3 -> nil
func PrintList(head *ListNode) {
	for cur := head; cur != nil; cur = cur.Next {
		fmt.Printf("%d -> ", cur.Val)
	}
	fmt.Println("nil")
}

// ============================================================
// BINARY TREE — TreeNode
// ============================================================
// LeetCode always gives you: type TreeNode struct { Val int; Left, Right *TreeNode }

// TreeNode represents a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// NewTree builds a binary tree from a level-order (BFS) slice.
// Use 0 to represent null nodes.
// Example: NewTree([]int{1, 2, 3, 0, 0, 4, 5})
//
//	  1
//	 / \
//	2   3
//	   / \
//	  4   5
func NewTree(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == 0 {
		return nil
	}
	root := &TreeNode{Val: vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if vals[i] != 0 {
				node.Left = &TreeNode{Val: vals[i]}
				queue = append(queue, node.Left)
			}
			i++
		}
		if i < len(vals) {
			if vals[i] != 0 {
				node.Right = &TreeNode{Val: vals[i]}
				queue = append(queue, node.Right)
			}
			i++
		}
	}
	return root
}

// TreeToLevelOrder returns the level-order values of a tree (BFS).
// Null nodes are represented as 0.
func TreeToLevelOrder(root *TreeNode) []int {
	if root == nil {
		return nil
	}
	var result []int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			result = append(result, 0)
		} else {
			result = append(result, node.Val)
			queue = append(queue, node.Left)
			queue = append(queue, node.Right)
		}
	}
	for len(result) > 0 && result[len(result)-1] == 0 {
		result = result[:len(result)-1]
	}
	return result
}

// ============================================================
// INTEGER HELPERS
// ============================================================

// Min returns the smaller of two ints.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two ints.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of n.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ============================================================
// SLICE / MATRIX HELPERS
// ============================================================

// PrintSlice prints a slice in a readable format.
func PrintSlice(s []int) {
	fmt.Print("[")
	for i, v := range s {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(v)
	}
	fmt.Println("]")
}

// PrintMatrix prints a 2D slice as a grid.
func PrintMatrix(m [][]int) {
	for _, row := range m {
		PrintSlice(row)
	}
}

// CopyMatrix returns a deep copy of a 2D int slice.
// Useful in tests when a function modifies the input in-place.
func CopyMatrix(m [][]int) [][]int {
	cp := make([][]int, len(m))
	for i, row := range m {
		cp[i] = make([]int, len(row))
		copy(cp[i], row)
	}
	return cp
}

// SlicesEqual returns true if two int slices have identical content.
func SlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
