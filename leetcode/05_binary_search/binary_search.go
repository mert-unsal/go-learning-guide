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

// ============================================================
// PROBLEM 6: Guess Number Higher or Lower (LeetCode #374) — EASY
// ============================================================
// We're playing a guessing game. You call guess(n) and it returns:
//   -1 if your guess is too high
//    1 if your guess is too low
//    0 if correct
// Find the number from 1 to n.

// guess is the API (simulated here via a closure approach in tests).
// For the solution, we define the function signature as LeetCode expects.

// GuessNumber returns the correct number using binary search.
// Time: O(log n)  Space: O(1)
func GuessNumber(n int, guessFn func(int) int) int {
	left, right := 1, n
	for left <= right {
		mid := left + (right-left)/2
		result := guessFn(mid)
		if result == 0 {
			return mid
		} else if result == -1 {
			right = mid - 1 // guess too high
		} else {
			left = mid + 1 // guess too low
		}
	}
	return -1
}

// ============================================================
// PROBLEM 7: Count Negative Numbers in a Sorted Matrix (LeetCode #1351) — EASY
// ============================================================
// Given a matrix sorted in non-increasing order (both rows and columns),
// count the number of negative numbers.
//
// Example: [[4,3,2,-1],[3,2,1,-1],[1,1,-1,-2],[-1,-1,-2,-3]] → 8
//
// Approach: start from top-right. If negative, all below are negative (count them).
// Move left. If positive, move down.

// CountNegatives counts negatives in a row+column sorted matrix.
// Time: O(m+n)  Space: O(1)
func CountNegatives(grid [][]int) int {
	rows, cols := len(grid), len(grid[0])
	r, c := 0, cols-1
	count := 0
	for r < rows && c >= 0 {
		if grid[r][c] < 0 {
			count += rows - r // all rows below are also negative in this column
			c--
		} else {
			r++
		}
	}
	return count
}

// ============================================================
// PROBLEM 8: Koko Eating Bananas (LeetCode #875) — MEDIUM
// ============================================================
// Koko has piles of bananas. She eats k bananas/hour. Find minimum k
// so she can eat all bananas within h hours.
//
// Example: piles=[3,6,7,11], h=8 → 4
//
// Key insight: binary search on the answer (eating speed k).
// For a given k, compute total hours needed; check if <= h.

// MinEatingSpeed returns the minimum eating speed to finish within h hours.
// Time: O(n log(max(piles)))  Space: O(1)
func MinEatingSpeed(piles []int, h int) int {
	maxPile := 0
	for _, p := range piles {
		if p > maxPile {
			maxPile = p
		}
	}

	canFinish := func(speed int) bool {
		hours := 0
		for _, pile := range piles {
			hours += (pile + speed - 1) / speed // ceiling division
		}
		return hours <= h
	}

	left, right := 1, maxPile
	for left < right {
		mid := left + (right-left)/2
		if canFinish(mid) {
			right = mid // mid might be the answer, keep it
		} else {
			left = mid + 1
		}
	}
	return left
}

// ============================================================
// PROBLEM 9: Find First and Last Position of Element (LeetCode #34) — MEDIUM
// ============================================================
// Given a sorted array, find the starting and ending position of target.
// Return [-1, -1] if not found.
//
// Example: nums=[5,7,7,8,8,10], target=8 → [3,4]
//
// Approach: binary search twice — once for left bound, once for right bound.

// SearchRange finds the first and last position of target in a sorted array.
// Time: O(log n)  Space: O(1)
func SearchRange(nums []int, target int) []int {
	return []int{findBound(nums, target, true), findBound(nums, target, false)}
}

func findBound(nums []int, target int, findLeft bool) int {
	left, right, result := 0, len(nums)-1, -1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			result = mid
			if findLeft {
				right = mid - 1 // keep searching left
			} else {
				left = mid + 1 // keep searching right
			}
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return result
}

// ============================================================
// PROBLEM 10: Median of Two Sorted Arrays (LeetCode #4) — HARD
// ============================================================
// Find the median of two sorted arrays. Must run in O(log(m+n)).
//
// Example: nums1=[1,3], nums2=[2] → 2.0
// Example: nums1=[1,2], nums2=[3,4] → 2.5
//
// Approach: binary search on the partition of the smaller array.
// We want a partition (i, j) such that:
//   max(left side) <= min(right side)
//   left side has (m+n+1)/2 elements
//
// Always binary search on the shorter array for efficiency.

// FindMedianSortedArrays returns the median of two sorted arrays.
// Time: O(log(min(m,n)))  Space: O(1)
func FindMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	// Ensure nums1 is the shorter array
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	m, n := len(nums1), len(nums2)
	half := (m + n + 1) / 2

	left, right := 0, m
	for left <= right {
		i := left + (right-left)/2 // partition in nums1
		j := half - i              // partition in nums2

		// Max of left side of nums1 partition
		maxLeft1 := -1 << 62
		if i > 0 {
			maxLeft1 = nums1[i-1]
		}
		// Min of right side of nums1 partition
		minRight1 := 1 << 62
		if i < m {
			minRight1 = nums1[i]
		}
		// Max of left side of nums2 partition
		maxLeft2 := -1 << 62
		if j > 0 {
			maxLeft2 = nums2[j-1]
		}
		// Min of right side of nums2 partition
		minRight2 := 1 << 62
		if j < n {
			minRight2 = nums2[j]
		}

		if maxLeft1 <= minRight2 && maxLeft2 <= minRight1 {
			// Correct partition found
			if (m+n)%2 == 1 {
				if maxLeft1 > maxLeft2 {
					return float64(maxLeft1)
				}
				return float64(maxLeft2)
			}
			maxLeft := maxLeft1
			if maxLeft2 > maxLeft {
				maxLeft = maxLeft2
			}
			minRight := minRight1
			if minRight2 < minRight {
				minRight = minRight2
			}
			return float64(maxLeft+minRight) / 2.0
		} else if maxLeft1 > minRight2 {
			right = i - 1 // move partition left
		} else {
			left = i + 1 // move partition right
		}
	}
	return 0.0
}
