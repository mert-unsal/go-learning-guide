// Package patterns contains reusable algorithm templates for common coding
// interview patterns. Study these scaffolds, then adapt them to specific problems.
//
// Run: go test ./patterns/... (no tests here — these are reference templates)
package patterns

import "container/heap"

// ============================================================
// 1. BINARY SEARCH TEMPLATE
// ============================================================
// Use when: sorted array, "find target / boundary / minimum satisfying condition"
//
// Key insight: always think about what mid represents and which half to discard.
// The `<=` in `left <= right` prevents missing the single-element case.

// BinarySearchExact finds target in a sorted slice. Returns index or -1.
// Time: O(log n)  Space: O(1)
func BinarySearchExact(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2 // avoids integer overflow vs (left+right)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1 // target is in right half
		} else {
			right = mid - 1 // target is in left half
		}
	}
	return -1 // not found
}

// BinarySearchLeftBound finds the leftmost index where nums[i] >= target.
// Returns len(nums) if all elements are less than target.
// Time: O(log n)  Space: O(1)
func BinarySearchLeftBound(nums []int, target int) int {
	left, right := 0, len(nums) // Note: right = len(nums), NOT len-1
	for left < right {          // Note: strict `<`, not `<=`
		mid := left + (right-left)/2
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid // keep mid as a candidate
		}
	}
	return left // left == right at this point
}

// ============================================================
// 2. SLIDING WINDOW TEMPLATE
// ============================================================
// Use when: subarray/substring with a constraint ("at most k", "all unique", etc.)
//
// Pattern: expand right pointer, shrink left pointer when constraint is violated.

// SlidingWindowMaxLen finds the length of the longest subarray satisfying
// a condition. Replace `isValid` with your actual constraint check.
// Time: O(n)  Space: O(1) to O(n) depending on what's tracked
func SlidingWindowMaxLen(nums []int) int {
	left := 0
	maxLen := 0
	// windowState could be a map, counter, sum, etc.
	windowSum := 0

	for right := 0; right < len(nums); right++ {
		// 1. Expand window: include nums[right]
		windowSum += nums[right]

		// 2. Shrink window from the left while constraint is violated
		for windowSum > 10 { // replace 10 with your constraint
			windowSum -= nums[left]
			left++
		}

		// 3. Update answer with current valid window size
		if right-left+1 > maxLen {
			maxLen = right - left + 1
		}
	}
	return maxLen
}

// ============================================================
// 3. TWO POINTERS TEMPLATE
// ============================================================
// Use when: sorted array, pair/triplet sum, palindrome check, merge

// TwoPointers demonstrates the converging-pointers pattern.
// Classic use: find pair with target sum in sorted array.
// Time: O(n)  Space: O(1)
func TwoPointers(nums []int, target int) (int, int) {
	left, right := 0, len(nums)-1
	for left < right {
		sum := nums[left] + nums[right]
		if sum == target {
			return left, right // found!
		} else if sum < target {
			left++ // need a larger sum → move left pointer right
		} else {
			right-- // need a smaller sum → move right pointer left
		}
	}
	return -1, -1 // not found
}

// ============================================================
// 4. BFS TEMPLATE
// ============================================================
// Use when: shortest path, level-order traversal, spreading in a grid
//
// Key: use a queue (slice as FIFO). Track visited to avoid cycles.

// BFSGrid performs BFS on a 2D grid from (startR, startC).
// Returns the shortest distance to any cell satisfying isGoal.
// Time: O(rows * cols)  Space: O(rows * cols)
func BFSGrid(grid [][]int, startR, startC int) int {
	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	type point struct{ r, c int }
	queue := []point{{startR, startC}}
	visited[startR][startC] = true
	dist := 0

	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} // right, left, down, up

	for len(queue) > 0 {
		// Process entire current level
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			cur := queue[0]
			queue = queue[1:]

			// Check if this is the goal
			if grid[cur.r][cur.c] == 9 { // replace 9 with your goal condition
				return dist
			}

			// Explore neighbors
			for _, d := range dirs {
				nr, nc := cur.r+d[0], cur.c+d[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
					visited[nr][nc] = true
					queue = append(queue, point{nr, nc})
				}
			}
		}
		dist++
	}
	return -1 // goal not reachable
}

// ============================================================
// 5. DFS TEMPLATE
// ============================================================
// Use when: explore all paths, connected components, cycle detection

// DFSRecursive performs DFS on a graph represented as adjacency list.
// Time: O(V + E)  Space: O(V) for the visited set + recursion stack
func DFSRecursive(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	var result []int

	var dfs func(node int)
	dfs = func(node int) {
		if visited[node] {
			return
		}
		visited[node] = true
		result = append(result, node)
		for _, neighbor := range graph[node] {
			dfs(neighbor)
		}
	}

	dfs(start)
	return result
}

// DFSIterative is the iterative DFS using an explicit stack.
// Useful when recursion depth might cause stack overflow.
// Time: O(V + E)  Space: O(V)
func DFSIterative(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	var result []int
	stack := []int{start}

	for len(stack) > 0 {
		// Pop from stack
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[node] {
			continue
		}
		visited[node] = true
		result = append(result, node)

		// Push neighbors (reverse order to match recursive DFS order)
		neighbors := graph[node]
		for i := len(neighbors) - 1; i >= 0; i-- {
			if !visited[neighbors[i]] {
				stack = append(stack, neighbors[i])
			}
		}
	}
	return result
}

// ============================================================
// 6. DYNAMIC PROGRAMMING TEMPLATE
// ============================================================
// Use when: overlapping subproblems + optimal substructure
//
// Two approaches:
//   Top-down (memoization): recursion + cache
//   Bottom-up (tabulation): iterative, fill dp table

// DPTopDown shows memoized recursion. Replace the recurrence with your own.
// Example: fibonacci(n) = fibonacci(n-1) + fibonacci(n-2)
// Time: O(n)  Space: O(n)
func DPTopDown(n int) int {
	memo := make(map[int]int)

	var solve func(i int) int
	solve = func(i int) int {
		if i <= 1 {
			return i // base case
		}
		if val, ok := memo[i]; ok {
			return val // return cached result
		}
		result := solve(i-1) + solve(i-2) // recurrence
		memo[i] = result                  // cache before returning
		return result
	}

	return solve(n)
}

// DPBottomUp shows tabulation. Fill the table from small to large.
// Example: fibonacci using 1D DP with O(1) space optimization.
// Time: O(n)  Space: O(1)
func DPBottomUp(n int) int {
	if n <= 1 {
		return n
	}
	prev2, prev1 := 0, 1
	for i := 2; i <= n; i++ {
		cur := prev1 + prev2 // recurrence
		prev2 = prev1
		prev1 = cur
	}
	return prev1
}

// ============================================================
// 7. MONOTONIC STACK TEMPLATE
// ============================================================
// Use when: "next greater element", "previous smaller element", histogram problems
//
// Invariant: stack stays sorted (increasing or decreasing).
// When we violate the invariant, we pop and process.

// NextGreaterElement finds the next greater element for each position.
// Returns -1 if no greater element exists to the right.
// Time: O(n)  Space: O(n)
func NextGreaterElement(nums []int) []int {
	n := len(nums)
	result := make([]int, n)
	for i := range result {
		result[i] = -1
	}
	stack := []int{} // stores INDICES (not values!)

	for i := 0; i < n; i++ {
		// While stack is non-empty and current element is greater
		// than the element at stack's top index → pop and record answer
		for len(stack) > 0 && nums[i] > nums[stack[len(stack)-1]] {
			idx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[idx] = nums[i] // nums[i] is the next greater for idx
		}
		stack = append(stack, i)
	}
	return result
}

// ============================================================
// 8. UNION-FIND (DISJOINT SET UNION) TEMPLATE
// ============================================================
// Use when: connected components, cycle detection in undirected graphs,
//           grouping/merging elements
//
// Optimizations: path compression + union by rank → near O(1) per operation

// UnionFind is a DSU data structure.
type UnionFind struct {
	parent []int
	rank   []int
	count  int // number of distinct components
}

// NewUnionFind initializes n elements each in their own component.
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i // each node is its own parent initially
	}
	return &UnionFind{parent: parent, rank: rank, count: n}
}

// Find returns the root of x's component (with path compression).
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // path compression
	}
	return uf.parent[x]
}

// Union merges the components of x and y. Returns true if they were separate.
func (uf *UnionFind) Union(x, y int) bool {
	rootX, rootY := uf.Find(x), uf.Find(y)
	if rootX == rootY {
		return false // already in the same component
	}
	// Union by rank: attach smaller tree under larger tree
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	uf.count--
	return true
}

// Connected returns true if x and y are in the same component.
func (uf *UnionFind) Connected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

// ============================================================
// 9. HEAP TEMPLATE
// ============================================================
// Use when: k-th largest/smallest, merge k sorted lists, priority queue
//
// Go's container/heap requires implementing the heap.Interface:
//   Len() int, Less(i,j int) bool, Swap(i,j int), Push(x any), Pop() any
//
// For a MAX heap, flip the Less comparison: return h[i] > h[j]

// MinHeap is a min-heap of integers (smallest element at top).
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] } // min at top
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// HeapExample shows how to use the MinHeap.
func HeapExample() int {
	h := &MinHeap{5, 3, 8, 1, 2}
	heap.Init(h) // heapify in O(n)

	heap.Push(h, 0) // push O(log n)

	min := heap.Pop(h).(int) // pop minimum O(log n)
	return min
}
