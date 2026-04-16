package graphs

// ============================================================
// PROBLEM 2: Course Schedule (LeetCode #207) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   There are numCourses courses labeled 0 to numCourses-1.
//   Prerequisites are given as pairs [a, b] meaning you must take
//   course b before course a. Return true if you can finish all
//   courses (i.e., the prerequisite graph has no cycles).
//
// PARAMETERS:
//   numCourses    int     — total number of courses
//   prerequisites [][]int — list of [course, prerequisite] pairs
//
// RETURN:
//   bool — true if all courses can be finished (no cycle in dependency graph)
//
// CONSTRAINTS:
//   • 1 <= numCourses <= 2000
//   • 0 <= prerequisites.length <= 5000
//   • prerequisites[i].length == 2
//   • 0 <= ai, bi < numCourses
//   • All prerequisite pairs are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  numCourses = 2, prerequisites = [[1,0]]
//   Output: true
//   Why:    Take course 0 first, then course 1. No cycle.
//
// Example 2:
//   Input:  numCourses = 2, prerequisites = [[1,0],[0,1]]
//   Output: false
//   Why:    Circular dependency: 0→1→0. Cannot finish.
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Model as directed graph: edge from prerequisite to course
// • DFS with 3-color marking: 0=unvisited, 1=in-current-path (cycle!), 2=fully-processed
// • Alternative: BFS topological sort (Kahn's algorithm) — if all nodes processed, no cycle
// • Target: O(V+E) time, O(V+E) space

func CanFinish(numCourses int, prerequisites [][]int) bool {
	return false
}
