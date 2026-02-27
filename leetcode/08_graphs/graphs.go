// Package graphs contains LeetCode graph problems with explanations.
// Topics: DFS flood-fill, BFS shortest path, topological sort, cycle detection.
package graphs

// ============================================================
// PROBLEM 1: Number of Islands (LeetCode #200) — MEDIUM
// ============================================================
// Given a 2D grid of '1' (land) and '0' (water), count the number of islands.
// An island is surrounded by water and formed by connecting adjacent lands
// horizontally or vertically.
//
// Example:
//   [["1","1","0"],
//    ["1","1","0"],
//    ["0","0","1"]]  → 2
//
// Approach: DFS flood-fill.
// For each unvisited '1', start a DFS that marks all connected '1's as visited
// (we overwrite with '0' in-place to avoid a separate visited array).
// Count how many DFS calls we make = number of islands.

// NumIslands counts the number of islands in the grid.
// Time: O(rows * cols)  Space: O(rows * cols) for recursion stack in worst case
func NumIslands(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	count := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '1' {
				count++
				dfsIsland(grid, r, c, rows, cols)
			}
		}
	}
	return count
}

// dfsIsland marks all connected land cells as visited (sets to '0').
func dfsIsland(grid [][]byte, r, c, rows, cols int) {
	// Boundary or water: stop
	if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] != '1' {
		return
	}
	grid[r][c] = '0' // mark visited
	dfsIsland(grid, r+1, c, rows, cols)
	dfsIsland(grid, r-1, c, rows, cols)
	dfsIsland(grid, r, c+1, rows, cols)
	dfsIsland(grid, r, c-1, rows, cols)
}

// ============================================================
// PROBLEM 2: Course Schedule (LeetCode #207) — MEDIUM
// ============================================================
// There are numCourses courses (0 to numCourses-1). prerequisites[i] = [a, b]
// means you must take b before a. Return true if you can finish all courses.
//
// This reduces to: does the directed graph have a cycle?
// If yes → impossible to finish. If no → topological order exists.
//
// Approach: DFS with 3-color marking.
//   0 = unvisited
//   1 = in current DFS path (gray) — if we reach gray, there's a cycle!
//   2 = fully processed (black) — safe, no cycle from here

// CanFinish returns true if all courses can be completed (no cycle in prereqs).
// Time: O(V + E)  Space: O(V + E)
func CanFinish(numCourses int, prerequisites [][]int) bool {
	// Build adjacency list
	graph := make([][]int, numCourses)
	for i := range graph {
		graph[i] = []int{}
	}
	for _, pre := range prerequisites {
		course, prereq := pre[0], pre[1]
		graph[prereq] = append(graph[prereq], course) // prereq → course
	}

	// 0=unvisited, 1=in-path (gray), 2=done (black)
	color := make([]int, numCourses)

	var hasCycle func(node int) bool
	hasCycle = func(node int) bool {
		if color[node] == 1 {
			return true // back edge → cycle!
		}
		if color[node] == 2 {
			return false // already fully explored, safe
		}
		color[node] = 1 // mark as in current path
		for _, neighbor := range graph[node] {
			if hasCycle(neighbor) {
				return true
			}
		}
		color[node] = 2 // fully processed
		return false
	}

	for i := 0; i < numCourses; i++ {
		if hasCycle(i) {
			return false
		}
	}
	return true
}

// ============================================================
// PROBLEM 3: Clone Graph (LeetCode #133) — MEDIUM
// ============================================================
// Given a reference node in an undirected connected graph, deep clone the graph.
// Each node has Val int and Neighbors []*Node.

// GraphNode is a graph node for the clone problem.
type GraphNode struct {
	Val       int
	Neighbors []*GraphNode
}

// CloneGraph returns a deep clone of the graph.
// Time: O(V + E)  Space: O(V) for the visited map
func CloneGraph(node *GraphNode) *GraphNode {
	if node == nil {
		return nil
	}
	// visited maps original node → cloned node
	visited := make(map[*GraphNode]*GraphNode)

	var clone func(n *GraphNode) *GraphNode
	clone = func(n *GraphNode) *GraphNode {
		if copy, ok := visited[n]; ok {
			return copy // already cloned
		}
		// Create new node
		cloned := &GraphNode{Val: n.Val}
		visited[n] = cloned // register BEFORE processing neighbors (handles cycles)

		for _, neighbor := range n.Neighbors {
			cloned.Neighbors = append(cloned.Neighbors, clone(neighbor))
		}
		return cloned
	}

	return clone(node)
}

// ============================================================
// PROBLEM 4: Flood Fill (LeetCode #733) — EASY
// ============================================================
// Given a 2D image, a starting pixel (sr, sc), and a new color,
// flood-fill the connected region (same original color as start).
//
// Approach: DFS from (sr, sc), change color of all connected same-color pixels.

// FloodFill performs flood fill from (sr, sc) with newColor.
// Time: O(rows * cols)  Space: O(rows * cols)
func FloodFill(image [][]int, sr int, sc int, color int) [][]int {
	originalColor := image[sr][sc]
	if originalColor == color {
		return image // already the target color, nothing to do
	}
	fill(image, sr, sc, originalColor, color, len(image), len(image[0]))
	return image
}

func fill(image [][]int, r, c, originalColor, newColor, rows, cols int) {
	if r < 0 || r >= rows || c < 0 || c >= cols || image[r][c] != originalColor {
		return
	}
	image[r][c] = newColor
	fill(image, r+1, c, originalColor, newColor, rows, cols)
	fill(image, r-1, c, originalColor, newColor, rows, cols)
	fill(image, r, c+1, originalColor, newColor, rows, cols)
	fill(image, r, c-1, originalColor, newColor, rows, cols)
}
