package graphs

// ============================================================
// PROBLEM 9: Course Schedule II (LeetCode #210) — MEDIUM
// ============================================================
//
// PROBLEM STATEMENT:
//   There are numCourses courses labeled 0 to numCourses-1. Given
//   prerequisite pairs [a, b] (must take b before a), return an
//   ordering of courses you should take to finish all courses. If
//   impossible (cycle exists), return an empty array. Any valid
//   topological order is accepted.
//
// PARAMETERS:
//   numCourses    int     — total number of courses
//   prerequisites [][]int — list of [course, prerequisite] pairs
//
// RETURN:
//   []int — valid course ordering, or empty slice if impossible
//
// CONSTRAINTS:
//   • 1 <= numCourses <= 2000
//   • 0 <= prerequisites.length <= numCourses * (numCourses - 1)
//   • prerequisites[i].length == 2
//   • 0 <= ai, bi < numCourses
//   • ai != bi
//   • All pairs are unique
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1:
//   Input:  numCourses = 2, prerequisites = [[1,0]]
//   Output: [0,1]
//   Why:    Must take 0 before 1
//
// Example 2:
//   Input:  numCourses = 4, prerequisites = [[1,0],[2,0],[3,1],[3,2]]
//   Output: [0,1,2,3] or [0,2,1,3]
//   Why:    Multiple valid orderings exist
//
// Example 3:
//   Input:  numCourses = 1, prerequisites = []
//   Output: [0]
//
// ─── HINTS ──────────────────────────────────────────────────
//
// • Topological sort via Kahn's algorithm: BFS with in-degree tracking
// • Or DFS post-order: nodes added in reverse finish order; detect cycles with 3-color marking
// • If result length < numCourses, a cycle exists → return empty
// • Target: O(V+E) time, O(V+E) space
func FindOrder(numCourses int, prerequisites [][]int) []int {
	return nil
}
