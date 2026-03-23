package arrays

// Merge ============================================================
// PROBLEM 6: Merge Sorted Array (LeetCode #88) — EASY
// ============================================================
//
// PROBLEM STATEMENT:
//
//	You are given two integer arrays nums1 and nums2, both sorted in
//	NON-DECREASING order, and two integers m and n representing the
//	number of valid elements in nums1 and nums2, respectively.
//
//	Merge nums2 into nums1 so that nums1 becomes a single sorted array.
//
//	The final sorted array must be stored INSIDE nums1.
//	To accommodate this, nums1 has a length of m + n, where:
//	  • The first m elements are the real data (sorted).
//	  • The last n elements are zeros — they are just placeholders and
//	    should be ignored. They exist solely to reserve space for merging.
//
//	You must NOT return a new array. The function has no return value.
//	The caller will inspect nums1 after your function returns.
//
// PARAMETERS:
//
//	nums1 []int — destination array, length = m + n.
//	               Positions [0..m-1] hold sorted data.
//	               Positions [m..m+n-1] hold placeholder zeros.
//	m     int   — how many real elements are in nums1.
//	nums2 []int — source array, length = n. All n elements are sorted.
//	n     int   — how many elements are in nums2.
//
// CONSTRAINTS:
//   - nums1.length == m + n
//   - nums2.length == n
//   - 0 <= m, n <= 200
//   - 1 <= m + n <= 200
//   - -10⁹ <= nums1[i], nums2[j] <= 10⁹
//
// WHY "IN-PLACE" MATTERS (Go-specific):
//
//	In Go, slices are passed by value — the function gets a COPY of the
//	slice header (pointer, length, capacity). If you create a new slice
//	and write `nums1 = newSlice`, you only reassign the local copy.
//	The caller's nums1 is unchanged. You MUST write into the existing
//	backing array of nums1 to make changes visible to the caller.
//
// ─── EXAMPLES ───────────────────────────────────────────────
//
// Example 1 — Basic merge:
//
//	Input:  nums1 = [1, 2, 3, 0, 0, 0],  m = 3
//	        nums2 = [2, 5, 6],            n = 3
//
//	What the array really looks like:
//	        real data → [1, 2, 3]  |  placeholders → [0, 0, 0]
//
//	After merge, nums1 should be: [1, 2, 2, 3, 5, 6]
//
// Example 2 — nums2 is empty (nothing to merge):
//
//	Input:  nums1 = [1],  m = 1
//	        nums2 = [],   n = 0
//	Result: nums1 = [1]   (unchanged)
//
// Example 3 — nums1 has no real data (just placeholders):
//
//	Input:  nums1 = [0],  m = 0
//	        nums2 = [1],  n = 1
//	Result: nums1 = [1]
//	Note: the zero in nums1 is a placeholder, not real data.
//
// Example 4 — All elements in nums2 are smaller than nums1:
//
//	Input:  nums1 = [4, 5, 6, 0, 0, 0],  m = 3
//	        nums2 = [1, 2, 3],            n = 3
//	Result: nums1 = [1, 2, 3, 4, 5, 6]
//	Every nums2 element must go before every nums1 element.
//
// Example 5 — Interleaved values across both arrays:
//
//	Input:  nums1 = [1, 3, 5, 7, 0, 0, 0, 0],  m = 4
//	        nums2 = [2, 4, 6, 8],                n = 4
//	Result: nums1 = [1, 2, 3, 4, 5, 6, 7, 8]
//	Elements alternate between the two sources.
//
// Example 6 — Duplicates across both arrays:
//
//	Input:  nums1 = [1, 2, 3, 0, 0, 0],  m = 3
//	        nums2 = [1, 2, 3],            n = 3
//	Result: nums1 = [1, 1, 2, 2, 3, 3]
//	Duplicates are perfectly valid; preserve all of them.
//
// Example 7 — Negative numbers and mixed signs:
//
//	Input:  nums1 = [-10, 0, 10, 0, 0, 0],  m = 3
//	        nums2 = [-5, 5, 15],             n = 3
//	Result: nums1 = [-10, -5, 0, 5, 10, 15]
//
// Example 8 — Large gap between values:
//
//	Input:  nums1 = [1, 1000000, 0, 0],  m = 2
//	        nums2 = [500, 999999],        n = 2
//	Result: nums1 = [1, 500, 999999, 1000000]
//
// Merge merges nums2 into nums1 in-place.
// Time: O(m+n)  Space: O(1)
func Merge(nums1 []int, m int, nums2 []int, n int) {
	i := m - 1
	j := len(nums2) - 1
	index := len(nums1) - 1
	for index >= 0 && j >= 0 {
		if i < 0 || nums2[j] >= nums1[i] {
			nums1[index] = nums2[j]
			j--
		} else {
			nums1[index] = nums1[i]
			i--
		}
		index--
	}
}
