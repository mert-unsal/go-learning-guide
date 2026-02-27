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

// ============================================================
// PROBLEM 5: Find if Path Exists in Graph (LeetCode #1971) — EASY
// ============================================================
// Given a bidirectional graph with n vertices and edges, determine if a
// valid path exists from source to destination.
//
// Example: n=3, edges=[[0,1],[1,2],[2,0]], source=0, destination=2 → true
//
// Approach: BFS/DFS from source, check if we reach destination.

// ValidPath returns true if a path exists from source to destination.
// Time: O(V + E)  Space: O(V + E)
func ValidPath(n int, edges [][]int, source int, destination int) bool {
	if source == destination {
		return true
	}
	// Build adjacency list
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	// BFS
	visited := make([]bool, n)
	queue := []int{source}
	visited[source] = true
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, next := range adj[curr] {
			if next == destination {
				return true
			}
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}

// ============================================================
// PROBLEM 6: Number of Connected Components in Undirected Graph (LeetCode #323) — MEDIUM
// ============================================================
// Given n nodes and a list of undirected edges, return the number of
// connected components.
//
// Example: n=5, edges=[[0,1],[1,2],[3,4]] → 2
//
// Approach: Union-Find (Disjoint Set Union).

// CountComponents returns the number of connected components.
// Time: O(n + e * α(n))  Space: O(n)
func CountComponents(n int, edges [][]int) int {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	components := n

	var find func(x int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x]) // path compression
		}
		return parent[x]
	}

	union := func(x, y int) {
		px, py := find(x), find(y)
		if px == py {
			return
		}
		components--
		if rank[px] < rank[py] {
			parent[px] = py
		} else if rank[px] > rank[py] {
			parent[py] = px
		} else {
			parent[py] = px
			rank[px]++
		}
	}

	for _, e := range edges {
		union(e[0], e[1])
	}
	return components
}

// ============================================================
// PROBLEM 7: Rotting Oranges (LeetCode #994) — MEDIUM
// ============================================================
// Grid of 0 (empty), 1 (fresh orange), 2 (rotten orange).
// Every minute, rotten oranges spread to adjacent fresh oranges.
// Return the minimum time until no fresh oranges remain, or -1 if impossible.
//
// Example: [[2,1,1],[1,1,0],[0,1,1]] → 4
//
// Approach: multi-source BFS. Start all rotten oranges in the queue simultaneously.

// OrangesRotting returns the minimum minutes to rot all oranges, or -1.
// Time: O(rows * cols)  Space: O(rows * cols)
func OrangesRotting(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	queue := [][2]int{}
	fresh := 0

	// Collect all initially rotten oranges and count fresh
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 2 {
				queue = append(queue, [2]int{r, c})
			} else if grid[r][c] == 1 {
				fresh++
			}
		}
	}

	if fresh == 0 {
		return 0
	}

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	minutes := 0

	for len(queue) > 0 && fresh > 0 {
		minutes++
		size := len(queue)
		for i := 0; i < size; i++ {
			curr := queue[i]
			for _, d := range dirs {
				nr, nc := curr[0]+d[0], curr[1]+d[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == 1 {
					grid[nr][nc] = 2
					fresh--
					queue = append(queue, [2]int{nr, nc})
				}
			}
		}
		queue = queue[size:]
	}

	if fresh > 0 {
		return -1
	}
	return minutes
}

// ============================================================
// PROBLEM 8: Pacific Atlantic Water Flow (LeetCode #417) — MEDIUM
// ============================================================
// An m×n island. Pacific ocean touches top/left edges, Atlantic touches bottom/right.
// Water flows to neighbors with equal or lower height.
// Return all cells from which water can flow to BOTH oceans.
//
// Approach: reverse BFS from both oceans. Mark cells reachable from Pacific,
// then from Atlantic. Return intersection.

// PacificAtlantic returns cells that can flow to both oceans.
// Time: O(m*n)  Space: O(m*n)
func PacificAtlantic(heights [][]int) [][]int {
	rows, cols := len(heights), len(heights[0])
	pacific := make([][]bool, rows)
	atlantic := make([][]bool, rows)
	for i := range pacific {
		pacific[i] = make([]bool, cols)
		atlantic[i] = make([]bool, cols)
	}

	bfs := func(queue [][2]int, visited [][]bool) {
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]
			for _, d := range dirs {
				nr, nc := curr[0]+d[0], curr[1]+d[1]
				if nr >= 0 && nr < rows && nc >= 0 && nc < cols &&
					!visited[nr][nc] && heights[nr][nc] >= heights[curr[0]][curr[1]] {
					visited[nr][nc] = true
					queue = append(queue, [2]int{nr, nc})
				}
			}
		}
	}

	// Seed Pacific (top row + left column)
	var pacQueue [][2]int
	var atlQueue [][2]int
	for r := 0; r < rows; r++ {
		pacific[r][0] = true
		atlantic[r][cols-1] = true
		pacQueue = append(pacQueue, [2]int{r, 0})
		atlQueue = append(atlQueue, [2]int{r, cols - 1})
	}
	for c := 0; c < cols; c++ {
		pacific[0][c] = true
		atlantic[rows-1][c] = true
		pacQueue = append(pacQueue, [2]int{0, c})
		atlQueue = append(atlQueue, [2]int{rows - 1, c})
	}

	bfs(pacQueue, pacific)
	bfs(atlQueue, atlantic)

	var result [][]int
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if pacific[r][c] && atlantic[r][c] {
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}
