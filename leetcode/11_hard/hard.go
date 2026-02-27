// Package hard contains LeetCode HARD level problems with detailed explanations.
// Topics: advanced DP, segment trees, hard graph problems, complex data structures.
package hard

import "sort"

// ============================================================
// PROBLEM 1: Merge k Sorted Lists (LeetCode #23) — HARD
// ============================================================
// Merge k sorted linked lists into one sorted list.
//
// Example: [[1,4,5],[1,3,4],[2,6]] → [1,1,2,3,4,4,5,6]
//
// Approach: divide and conquer — merge pairs of lists repeatedly.
// Reduces to O(n log k) where n = total nodes, k = number of lists.
// Merging two sorted lists is O(n+m); we do log k rounds.

// ListNode is a singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// MergeKLists merges k sorted linked lists into one sorted list.
// Time: O(n log k)  Space: O(log k) recursion
func MergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	return mergeRange(lists, 0, len(lists)-1)
}

func mergeRange(lists []*ListNode, left, right int) *ListNode {
	if left == right {
		return lists[left]
	}
	mid := left + (right-left)/2
	l1 := mergeRange(lists, left, mid)
	l2 := mergeRange(lists, mid+1, right)
	return mergeTwoLists(l1, l2)
}

func mergeTwoLists(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			cur.Next = l1
			l1 = l1.Next
		} else {
			cur.Next = l2
			l2 = l2.Next
		}
		cur = cur.Next
	}
	if l1 != nil {
		cur.Next = l1
	} else {
		cur.Next = l2
	}
	return dummy.Next
}

// ============================================================
// PROBLEM 2: Trapping Rain Water (LeetCode #42) — HARD
// ============================================================
// Given n non-negative integers representing an elevation map,
// compute how much water can be trapped after raining.
//
// Example: height=[0,1,0,2,1,0,1,3,2,1,2,1] → 6
//
// Approach: two-pointer O(1) space.
// Water at position i = min(maxLeft[i], maxRight[i]) - height[i]
// Two pointers from both ends. Process the side with the smaller max boundary.

// Trap returns the total units of trapped rain water.
// Time: O(n)  Space: O(1)
func Trap(height []int) int {
	if len(height) == 0 {
		return 0
	}
	left, right := 0, len(height)-1
	maxLeft, maxRight := 0, 0
	water := 0

	for left < right {
		if height[left] < height[right] {
			if height[left] >= maxLeft {
				maxLeft = height[left]
			} else {
				water += maxLeft - height[left]
			}
			left++
		} else {
			if height[right] >= maxRight {
				maxRight = height[right]
			} else {
				water += maxRight - height[right]
			}
			right--
		}
	}
	return water
}

// ============================================================
// PROBLEM 3: Word Ladder (LeetCode #127) — HARD
// ============================================================
// Given a beginWord, endWord, and wordList, find the shortest transformation
// sequence length from beginWord to endWord. Each step: change one letter,
// and the result must be in wordList.
//
// Example: beginWord="hit", endWord="cog",
//          wordList=["hot","dot","dog","lot","log","cog"] → 5
//          (hit→hot→dot→dog→cog)
//
// Approach: BFS — each level transforms one character.
// Use a word set for O(1) lookup. Remove words when visited.

// LadderLength returns the shortest transformation sequence length, or 0.
// Time: O(n * m²) where n = wordList size, m = word length  Space: O(n*m)
func LadderLength(beginWord string, endWord string, wordList []string) int {
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}
	if !wordSet[endWord] {
		return 0
	}

	queue := []string{beginWord}
	steps := 1

	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			word := queue[i]
			wordBytes := []byte(word)
			for j := 0; j < len(wordBytes); j++ {
				original := wordBytes[j]
				for c := byte('a'); c <= 'z'; c++ {
					if c == original {
						continue
					}
					wordBytes[j] = c
					next := string(wordBytes)
					if next == endWord {
						return steps + 1
					}
					if wordSet[next] {
						queue = append(queue, next)
						delete(wordSet, next) // mark visited
					}
					wordBytes[j] = original // restore
				}
			}
		}
		queue = queue[size:]
		steps++
	}
	return 0
}

// ============================================================
// PROBLEM 4: Longest Valid Parentheses (LeetCode #32) — HARD
// ============================================================
// Find the length of the longest valid parentheses substring.
//
// Example: s=")()())" → 4  ("()()")
// Example: s="(()" → 2   ("()")
//
// Approach: stack-based.
// Push indices onto the stack. Maintain a "base" index for valid substring calculation.
// When we see '(', push its index.
// When we see ')':
//   - if stack is non-empty, pop → length = i - stack.top
//   - if stack is empty, push i as new base

// LongestValidParentheses returns the length of the longest valid substring.
// Time: O(n)  Space: O(n)
func LongestValidParentheses(s string) int {
	stack := []int{-1} // base index
	maxLen := 0

	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			stack = append(stack, i)
		} else {
			stack = stack[:len(stack)-1] // pop
			if len(stack) == 0 {
				stack = append(stack, i) // new base
			} else {
				length := i - stack[len(stack)-1]
				if length > maxLen {
					maxLen = length
				}
			}
		}
	}
	return maxLen
}

// ============================================================
// PROBLEM 5: Jump Game II (LeetCode #45) — MEDIUM/HARD
// ============================================================
// Find the minimum number of jumps to reach the last index.
// You can always reach the last index.
//
// Example: nums=[2,3,1,1,4] → 2  (jump 2 → jump 4)
//
// Greedy: track the current jump's reachable boundary and the next jump's farthest.
// Increment jumps when we reach the current boundary.

// Jump returns the minimum number of jumps to reach the last index.
// Time: O(n)  Space: O(1)
func Jump(nums []int) int {
	jumps, currentEnd, farthest := 0, 0, 0
	for i := 0; i < len(nums)-1; i++ {
		if i+nums[i] > farthest {
			farthest = i + nums[i]
		}
		if i == currentEnd {
			jumps++
			currentEnd = farthest
			if currentEnd >= len(nums)-1 {
				break
			}
		}
	}
	return jumps
}

// ============================================================
// PROBLEM 6: N-Queens (LeetCode #51) — HARD
// ============================================================
// Place n queens on an n×n board so no two queens attack each other.
// Return all distinct solutions. Each solution is represented as a slice
// of strings where 'Q' = queen and '.' = empty.
//
// Example: n=4 →
//   [".Q..","...Q","Q...","..Q."]
//   ["..Q.","Q...","...Q",".Q.."]
//
// Approach: backtracking. Track which columns, diagonals, anti-diagonals are attacked.

// SolveNQueens returns all solutions to the N-Queens problem.
// Time: O(n!)  Space: O(n²)
func SolveNQueens(n int) [][]string {
	var result [][]string
	cols := make([]bool, n)
	diag1 := make([]bool, 2*n) // row - col + n
	diag2 := make([]bool, 2*n) // row + col
	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	var backtrack func(row int)
	backtrack = func(row int) {
		if row == n {
			// Found a valid placement — convert to strings
			sol := make([]string, n)
			for i, r := range board {
				sol[i] = string(r)
			}
			result = append(result, sol)
			return
		}
		for col := 0; col < n; col++ {
			d1, d2 := row-col+n, row+col
			if cols[col] || diag1[d1] || diag2[d2] {
				continue
			}
			// Place queen
			board[row][col] = 'Q'
			cols[col], diag1[d1], diag2[d2] = true, true, true
			backtrack(row + 1)
			// Remove queen
			board[row][col] = '.'
			cols[col], diag1[d1], diag2[d2] = false, false, false
		}
	}
	backtrack(0)
	return result
}

// ============================================================
// PROBLEM 7: Serialize and Deserialize Binary Tree (LeetCode #297) — HARD
// ============================================================
// Design an algorithm to serialize a binary tree to a string and
// deserialize it back.
//
// Approach: preorder DFS with null markers.
// Serialize: "1,2,null,null,3,null,null"
// Deserialize: split by comma, build recursively.

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Serialize converts a binary tree to a string.
// Time: O(n)  Space: O(n)
func Serialize(root *TreeNode) string {
	var result []byte
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			result = append(result, []byte("null,")...)
			return
		}
		// Append val + comma
		val := node.Val
		if val < 0 {
			result = append(result, '-')
			val = -val
		}
		digits := []byte{}
		if val == 0 {
			digits = []byte{'0'}
		}
		for val > 0 {
			digits = append([]byte{byte('0' + val%10)}, digits...)
			val /= 10
		}
		result = append(result, digits...)
		result = append(result, ',')
		dfs(node.Left)
		dfs(node.Right)
	}
	dfs(root)
	if len(result) > 0 {
		result = result[:len(result)-1] // remove trailing comma
	}
	return string(result)
}

// Deserialize converts a serialized string back to a binary tree.
// Time: O(n)  Space: O(n)
func Deserialize(data string) *TreeNode {
	tokens := splitByComma(data)
	idx := 0
	var build func() *TreeNode
	build = func() *TreeNode {
		if idx >= len(tokens) || tokens[idx] == "null" {
			idx++
			return nil
		}
		val := parseIntSimple(tokens[idx])
		idx++
		node := &TreeNode{Val: val}
		node.Left = build()
		node.Right = build()
		return node
	}
	return build()
}

func splitByComma(s string) []string {
	var tokens []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			tokens = append(tokens, s[start:i])
			start = i + 1
		}
	}
	return tokens
}

func parseIntSimple(s string) int {
	neg := false
	start := 0
	if len(s) > 0 && s[0] == '-' {
		neg = true
		start = 1
	}
	val := 0
	for _, ch := range s[start:] {
		val = val*10 + int(ch-'0')
	}
	if neg {
		return -val
	}
	return val
}

// ============================================================
// PROBLEM 8: Minimum Window Substring (LeetCode #76) — HARD
// ============================================================
// (Already in sliding_window — canonical hard problem, restated here)
// Find the smallest window in s containing all characters of t.
//
// Example: s="ADOBECODEBANC", t="ABC" → "BANC"
//
// Approach: sliding window with character frequency maps.

// MinWindow returns the minimum window substring containing all chars of t.
// Time: O(|s| + |t|)  Space: O(|s| + |t|)
func MinWindow(s string, t string) string {
	if len(s) == 0 || len(t) == 0 {
		return ""
	}
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	required := len(need)
	have := make(map[byte]int)
	formed := 0
	left := 0
	minLen := len(s) + 1
	minLeft := 0

	for right := 0; right < len(s); right++ {
		ch := s[right]
		have[ch]++
		if cnt, ok := need[ch]; ok && have[ch] == cnt {
			formed++
		}
		for formed == required {
			if right-left+1 < minLen {
				minLen = right - left + 1
				minLeft = left
			}
			lc := s[left]
			have[lc]--
			if cnt, ok := need[lc]; ok && have[lc] < cnt {
				formed--
			}
			left++
		}
	}
	if minLen == len(s)+1 {
		return ""
	}
	return s[minLeft : minLeft+minLen]
}

// ============================================================
// PROBLEM 9: Alien Dictionary (LeetCode #269) — HARD
// ============================================================
// Given a sorted list of words in an alien language, derive the character order.
// Return any valid ordering, or "" if invalid.
//
// Example: words=["wrt","wrf","er","ett","rftt"] → "wertf"
//
// Approach: build a directed graph from adjacent word pairs, then topological sort.
// If chars differ at position i, we know words[i][i] → words[i+1][i] ordering.
// Cycle = invalid (return "").

// AlienOrder returns a valid character ordering, or "" if impossible.
// Time: O(C) where C = total characters  Space: O(1) — at most 26 chars
func AlienOrder(words []string) string {
	// Initialize in-degree and adjacency list for all unique chars
	inDegree := make(map[byte]int)
	adj := make(map[byte][]byte)
	for _, word := range words {
		for i := 0; i < len(word); i++ {
			if _, exists := inDegree[word[i]]; !exists {
				inDegree[word[i]] = 0
				adj[word[i]] = []byte{}
			}
		}
	}

	// Build edges from adjacent words
	for i := 0; i < len(words)-1; i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		// Invalid: longer word is prefix of shorter word
		if len(w1) > len(w2) && w1[:minLen] == w2[:minLen] {
			return ""
		}
		for j := 0; j < minLen; j++ {
			if w1[j] != w2[j] {
				adj[w1[j]] = append(adj[w1[j]], w2[j])
				inDegree[w2[j]]++
				break
			}
		}
	}

	// BFS topological sort (Kahn's algorithm)
	queue := []byte{}
	for ch, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, ch)
		}
	}
	// Sort for deterministic output
	sort.Slice(queue, func(i, j int) bool { return queue[i] < queue[j] })

	result := []byte{}
	for len(queue) > 0 {
		ch := queue[0]
		queue = queue[1:]
		result = append(result, ch)
		neighbors := adj[ch]
		sort.Slice(neighbors, func(i, j int) bool { return neighbors[i] < neighbors[j] })
		for _, next := range neighbors {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(result) != len(inDegree) {
		return "" // cycle detected
	}
	return string(result)
}

// ============================================================
// PROBLEM 10: Regular Expression Matching (LeetCode #10) — HARD
// ============================================================
// Implement regex matching with '.' (any single char) and '*' (zero or more
// of the preceding element).
//
// Example: s="aa", p="a*" → true
// Example: s="ab", p=".*" → true
// Example: s="aab", p="c*a*b" → true
//
// Approach: 2D DP.
// dp[i][j] = true if s[0..i-1] matches p[0..j-1]
// If p[j-1] == '*':
//   Case 1: use zero occurrences of p[j-2]: dp[i][j] = dp[i][j-2]
//   Case 2: use p[j-2] once more: dp[i][j] |= dp[i-1][j] if p[j-2]=='.' or p[j-2]==s[i-1]
// Else if p[j-1]=='.' or p[j-1]==s[i-1]: dp[i][j] = dp[i-1][j-1]

// IsMatch returns true if s fully matches pattern p.
// Time: O(m*n)  Space: O(m*n)
func IsMatch(s string, p string) bool {
	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true // empty matches empty

	// Patterns like a*, a*b*, a*b*c* can match empty string
	for j := 2; j <= n; j++ {
		if p[j-1] == '*' {
			dp[0][j] = dp[0][j-2]
		}
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == '*' {
				// Zero occurrences of p[j-2]
				dp[i][j] = dp[i][j-2]
				// One more occurrence of p[j-2]
				if p[j-2] == '.' || p[j-2] == s[i-1] {
					dp[i][j] = dp[i][j] || dp[i-1][j]
				}
			} else if p[j-1] == '.' || p[j-1] == s[i-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}
	return dp[m][n]
}

// ============================================================
// PROBLEM 11: Edit Distance (LeetCode #72) — HARD
// ============================================================
// Find the minimum number of operations (insert, delete, replace) to
// convert word1 to word2.
//
// Example: word1="horse", word2="ros" → 3
//
// dp[i][j] = min edit distance between word1[0..i-1] and word2[0..j-1]
// If chars match: dp[i][j] = dp[i-1][j-1]
// Else: dp[i][j] = 1 + min(dp[i-1][j],   // delete from word1
//                          dp[i][j-1],     // insert into word1
//                          dp[i-1][j-1])   // replace

// MinDistance returns the minimum edit distance between word1 and word2.
// Time: O(m*n)  Space: O(m*n)
func MinDistance(word1 string, word2 string) int {
	m, n := len(word1), len(word2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	// Base cases: converting to/from empty string
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1] // characters match, no op
			} else {
				dp[i][j] = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return dp[m][n]
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

// ============================================================
// PROBLEM 12: Maximum Profit in Job Scheduling (LeetCode #1235) — HARD
// ============================================================
// Given jobs with startTime, endTime, and profit, find the maximum profit
// you can achieve without scheduling two overlapping jobs.
//
// Example: startTime=[1,2,3,3], endTime=[3,4,5,6], profit=[50,10,40,70] → 120
//
// Approach: sort by endTime + DP + binary search.
// dp[i] = max profit considering first i jobs (sorted by end time).
// For each job i, either skip it (dp[i-1]) or take it:
//   find the latest job j that ends <= startTime[i], profit = dp[j] + profit[i]

// JobScheduling returns the maximum profit from non-overlapping jobs.
// Time: O(n log n)  Space: O(n)
func JobScheduling(startTime []int, endTime []int, profit []int) int {
	n := len(startTime)
	type Job struct{ start, end, profit int }
	jobs := make([]Job, n)
	for i := range jobs {
		jobs[i] = Job{startTime[i], endTime[i], profit[i]}
	}
	sort.Slice(jobs, func(i, j int) bool { return jobs[i].end < jobs[j].end })

	// dp[i] = max profit using first i jobs
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		job := jobs[i-1]
		// Binary search: find latest job ending <= job.start
		lo, hi := 0, i-1
		for lo < hi {
			mid := lo + (hi-lo+1)/2
			if jobs[mid-1].end <= job.start {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		take := dp[lo] + job.profit
		skip := dp[i-1]
		if take > skip {
			dp[i] = take
		} else {
			dp[i] = skip
		}
	}
	return dp[n]
}
