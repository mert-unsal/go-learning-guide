package hard

import "sort"

var _ = sort.Ints

// Shared types.
type ListNode struct {
	Val  int
	Next *ListNode
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// PROBLEM 1: Merge k Sorted Lists (LeetCode #23) — HARD
func MergeKLists(lists []*ListNode) *ListNode                 { return nil }
func mergeRange(lists []*ListNode, left, right int) *ListNode { return nil }
func mergeTwoLists(l1, l2 *ListNode) *ListNode                { return nil }

// PROBLEM 2: Trapping Rain Water (LeetCode #42) — HARD
func Trap(height []int) int { return 0 }

// PROBLEM 3: Word Ladder (LeetCode #127) — HARD
func LadderLength(beginWord string, endWord string, wordList []string) int { return 0 }

// PROBLEM 4: Longest Valid Parentheses (LeetCode #32) — HARD
func LongestValidParentheses(s string) int { return 0 }

// PROBLEM 5: Jump Game II (LeetCode #45) — MEDIUM/HARD
func Jump(nums []int) int { return 0 }

// PROBLEM 6: N-Queens (LeetCode #51) — HARD
func SolveNQueens(n int) [][]string { return nil }

// PROBLEM 7: Serialize and Deserialize Binary Tree (LeetCode #297) — HARD
func Serialize(root *TreeNode) string   { return "" }
func Deserialize(data string) *TreeNode { return nil }
func splitByComma(s string) []string    { return nil }
func parseIntSimple(s string) int       { return 0 }

// PROBLEM 8: Minimum Window Substring (LeetCode #76) — HARD
func MinWindow(s string, t string) string { return "" }

// PROBLEM 9: Alien Dictionary (LeetCode #269) — HARD
func AlienOrder(words []string) string { return "" }

// PROBLEM 10: Regular Expression Matching (LeetCode #10) — HARD
func IsMatch(s string, p string) bool { return false }

// PROBLEM 11: Edit Distance (LeetCode #72) — HARD
func MinDistance(word1 string, word2 string) int { return 0 }
func min3(a, b, c int) int                       { return 0 }

// PROBLEM 12: Maximum Profit in Job Scheduling (LeetCode #1235) — HARD
func JobScheduling(startTime []int, endTime []int, profit []int) int { return 0 }
