// Package backtracking contains LeetCode backtracking problems.
// Topics: subsets, permutations, combinations, constraint satisfaction.
package backtracking

import "sort"

// Suppress unused import — you will need sort for some problems.
var _ = sort.Ints

// ============================================================
// PROBLEM 1: Subsets (LeetCode #78) — MEDIUM
// ============================================================
// Given an integer array of unique elements, return all possible subsets (power set).
//
// Example: nums=[1,2,3] → [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
//
// Approach: backtracking. At each position, choose to include or exclude the element.

// Subsets returns all possible subsets.
// Time: O(n * 2^n)  Space: O(n) recursion depth
func Subsets(nums []int) [][]int {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 2: Permutations (LeetCode #46) — MEDIUM
// ============================================================
// Given an array of distinct integers, return all possible permutations.
//
// Example: nums=[1,2,3] → [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
//
// Approach: backtracking with a used-set.

// Permute returns all permutations of nums.
// Time: O(n * n!)  Space: O(n)
func Permute(nums []int) [][]int {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 3: Combination Sum (LeetCode #39) — MEDIUM
// ============================================================
// Given candidates (no duplicates) and a target, find all unique combinations
// that sum to target. Each candidate can be used unlimited times.
//
// Example: candidates=[2,3,6,7], target=7 → [[2,2,3],[7]]
//
// Approach: backtracking. At each step, include candidate[i] (may reuse) or skip.

// CombinationSum returns all unique combinations summing to target.
// Time: O(n^(target/min))  Space: O(target/min)
func CombinationSum(candidates []int, target int) [][]int {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 4: Combination Sum II (LeetCode #40) — MEDIUM
// ============================================================
// Like #39 but candidates may contain duplicates and each number
// can be used at most once. Find all unique combinations summing to target.
//
// Example: candidates=[10,1,2,7,6,1,5], target=8 → [[1,1,6],[1,2,5],[1,7],[2,6]]

// CombinationSum2 returns all unique combinations (no reuse) summing to target.
// Time: O(2^n)  Space: O(n)
func CombinationSum2(candidates []int, target int) [][]int {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 5: Word Search (LeetCode #79) — MEDIUM
// ============================================================
// Given a 2D board of characters and a word, return true if the word
// exists in the grid. The word must be formed from adjacent cells
// (horizontally or vertically). Each cell may be used at most once.
//
// Example: board=[["A","B","C","E"],["S","F","C","S"],["A","D","E","E"]], word="ABCCED" → true

// Exist returns true if word exists in the board.
// Time: O(m*n * 3^L) where L = word length  Space: O(L)
func Exist(board [][]byte, word string) bool {
	// TODO: implement
	return false
}

// ============================================================
// PROBLEM 6: Letter Combinations of a Phone Number (LeetCode #17) — MEDIUM
// ============================================================
// Given a string containing digits 2-9, return all letter combinations
// that the digits could represent.
//
// Example: digits="23" → ["ad","ae","af","bd","be","bf","cd","ce","cf"]

// LetterCombinations returns all letter combinations for the given digits.
// Time: O(4^n)  Space: O(n)
func LetterCombinations(digits string) []string {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 7: Palindrome Partitioning (LeetCode #131) — MEDIUM
// ============================================================
// Partition a string such that every substring is a palindrome.
// Return all possible palindrome partitionings.
//
// Example: s="aab" → [["a","a","b"],["aa","b"]]

// Partition returns all palindrome partitionings of s.
// Time: O(n * 2^n)  Space: O(n)
func Partition(s string) [][]string {
	// TODO: implement
	return nil
}

// ============================================================
// PROBLEM 8: Subsets II (LeetCode #90) — MEDIUM
// ============================================================
// Given an integer array that may contain duplicates, return all possible
// subsets. The solution set must not contain duplicate subsets.
//
// Example: nums=[1,2,2] → [[],[1],[1,2],[1,2,2],[2],[2,2]]

// SubsetsWithDup returns all unique subsets of nums (may contain duplicates).
// Time: O(n * 2^n)  Space: O(n)
func SubsetsWithDup(nums []int) [][]int {
	// TODO: implement
	return nil
}
