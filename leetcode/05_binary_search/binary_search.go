// Package binary_search contains LeetCode binary search problems.
// Topics: rotated array search, boundary finding, search space reduction.
package binary_search

// ============================================================
// PROBLEM 1: Search in Rotated Sorted Array (LeetCode #33) — MEDIUM
// ============================================================
// A sorted array was rotated at some pivot. Find target. Return index or -1.
//
// Example: nums=[4,5,6,7,0,1,2], target=0 → 4
// Example: nums=[4,5,6,7,0,1,2], target=3 → -1
//
// Key insight: one of the two halves is ALWAYS sorted.
// Identify which half is sorted, then check if target is in that half.
// Narrow search space accordingly.

// Search finds target in a rotated sorted array. Returns index or -1.
// Time: O(log n)  Space: O(1)
func Search(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			return mid
		}

		// Determine which half is sorted
		if nums[left] <= nums[mid] {
			// Left half [left..mid] is sorted
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1 // target is in the sorted left half
			} else {
				left = mid + 1 // target is in the right half
			}
		} else {
			// Right half [mid..right] is sorted
			if nums[mid] < target && target <= nums[right] {
				left = mid + 1 // target is in the sorted right half
			} else {
				right = mid - 1 // target is in the left half
			}
		}
	}
	return -1
}

// ============================================================
// PROBLEM 2: Find Minimum in Rotated Sorted Array (LeetCode #153) — MEDIUM
// ============================================================
// Find the minimum element in a rotated sorted array (no duplicates).
//
// Example: nums=[3,4,5,1,2] → 1
// Example: nums=[4,5,6,7,0,1,2] → 0
// Example: nums=[11,13,15,17] → 11 (not rotated)
//
// Key insight: the minimum is at the "rotation point".
// If nums[mid] > nums[right], the minimum is in the right half.
// Otherwise, the minimum is in the left half (including mid).

// FindMin returns the minimum element in a rotated sorted array.
// Time: O(log n)  Space: O(1)
func FindMin(nums []int) int {
	left, right := 0, len(nums)-1

	for left < right {
		mid := left + (right-left)/2

		if nums[mid] > nums[right] {
			// Minimum is in the right half (mid cannot be minimum)
			left = mid + 1
		} else {
			// Minimum is in the left half (mid could be minimum)
			right = mid
		}
	}
	return nums[left] // left == right == index of minimum
}

// ============================================================
// PROBLEM 3: Binary Search (LeetCode #704) — EASY
// ============================================================
// Classic binary search on a sorted array.

// BinarySearch returns the index of target in sorted nums, or -1.
// Time: O(log n)  Space: O(1)
func BinarySearch(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// ============================================================
// PROBLEM 4: Search a 2D Matrix (LeetCode #74) — MEDIUM
// ============================================================
// m×n matrix where each row is sorted and each row's first element
// is greater than the previous row's last. Search for target.
//
// Approach: treat the matrix as a single sorted array of m*n elements.
// Map index i → row: i/n, col: i%n

// SearchMatrix returns true if target exists in the sorted matrix.
// Time: O(log(m*n))  Space: O(1)
func SearchMatrix(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}
	m, n := len(matrix), len(matrix[0])
	left, right := 0, m*n-1

	for left <= right {
		mid := left + (right-left)/2
		val := matrix[mid/n][mid%n] // convert 1D index to 2D
		if val == target {
			return true
		} else if val < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return false
}

// ============================================================
// PROBLEM 5: Find Peak Element (LeetCode #162) — MEDIUM
// ============================================================
// A peak element is greater than its neighbors. Find any peak index.
// nums[-1] and nums[n] are treated as -∞.
//
// Approach: binary search. If nums[mid] < nums[mid+1], the peak is
// to the right (we're on an ascending slope). Otherwise it's left or mid.

// FindPeakElement returns the index of any peak element.
// Time: O(log n)  Space: O(1)
func FindPeakElement(nums []int) int {
	left, right := 0, len(nums)-1
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] < nums[mid+1] {
			left = mid + 1 // ascending slope, peak is to the right
		} else {
			right = mid // descending slope, peak is at mid or to the left
		}
	}
	return left
}
