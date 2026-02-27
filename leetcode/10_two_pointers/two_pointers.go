// Package two_pointers contains LeetCode two-pointer problems.
// Topics: converging pointers, sort+two-pointer, water trapping.
package two_pointers

import "sort"

// ============================================================
// PROBLEM 1: Container With Most Water (LeetCode #11) — MEDIUM
// ============================================================
// Given an array of heights, find two lines that together with the x-axis
// form a container that holds the most water.
//
// Example: height=[1,8,6,2,5,4,8,3,7] → 49
//
// Approach: two pointers starting at both ends.
// Area = min(height[left], height[right]) * (right - left)
// Move the pointer with the SHORTER height inward (keeping the taller
// might give a larger area; the shorter can never benefit from staying).

// MaxArea returns the maximum water a container can hold.
// Time: O(n)  Space: O(1)
func MaxArea(height []int) int {
	left, right := 0, len(height)-1
	maxWater := 0

	for left < right {
		h := height[left]
		if height[right] < h {
			h = height[right]
		}
		area := h * (right - left)
		if area > maxWater {
			maxWater = area
		}
		// Move the shorter wall inward
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return maxWater
}

// ============================================================
// PROBLEM 2: 3Sum (LeetCode #15) — MEDIUM
// ============================================================
// Find all unique triplets in nums that sum to zero.
//
// Example: nums=[-1,0,1,2,-1,-4] → [[-1,-1,2],[-1,0,1]]
//
// Approach: sort + two-pointer.
// Fix the first element (iterate i from 0 to n-3).
// Use two pointers left=i+1, right=n-1 to find pairs summing to -nums[i].
// Skip duplicate values to avoid duplicate triplets.

// ThreeSum returns all unique triplets that sum to zero.
// Time: O(n²)  Space: O(1) extra (excluding output)
func ThreeSum(nums []int) [][]int {
	sort.Ints(nums) // sort first!
	var result [][]int

	for i := 0; i < len(nums)-2; i++ {
		// Skip duplicate first elements
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		left, right := i+1, len(nums)-1

		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum == 0 {
				result = append(result, []int{nums[i], nums[left], nums[right]})
				// Skip duplicates for second and third elements
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}
				left++
				right--
			} else if sum < 0 {
				left++ // need a larger sum
			} else {
				right-- // need a smaller sum
			}
		}
	}
	return result
}

// ============================================================
// PROBLEM 3: Trapping Rain Water (LeetCode #42) — HARD
// ============================================================
// Given n non-negative integers representing an elevation map,
// compute how much water can be trapped after raining.
//
// Example: height=[0,1,0,2,1,0,1,3,2,1,2,1] → 6
//
// Approach: two-pointer O(1) space.
// Water at position i = min(maxLeft[i], maxRight[i]) - height[i]
//
// Two pointers l=0, r=n-1. Track maxLeft and maxRight seen so far.
// If maxLeft < maxRight, the water at l is determined by maxLeft:
//   water += maxLeft - height[l], then advance l.
// Otherwise water at r is determined by maxRight: advance r.

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
				maxLeft = height[left] // update running max
			} else {
				water += maxLeft - height[left] // water fills up to maxLeft
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
// PROBLEM 4: Move Zeroes (LeetCode #283) — EASY
// ============================================================
// Move all 0's to the end while maintaining the relative order of non-zero elements.
// Do it in-place.
//
// Approach: slow pointer tracks the position for the next non-zero element.

// MoveZeroes moves all zeros to the end in-place.
// Time: O(n)  Space: O(1)
func MoveZeroes(nums []int) {
	slow := 0 // position for next non-zero
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != 0 {
			nums[slow] = nums[fast]
			slow++
		}
	}
	// Fill remaining positions with zeros
	for slow < len(nums) {
		nums[slow] = 0
		slow++
	}
}

// ============================================================
// PROBLEM 5: Remove Duplicates from Sorted Array (LeetCode #26) — EASY
// ============================================================
// Remove duplicates in-place from a sorted array. Return the new length.
//
// Approach: slow pointer is the position for the next unique element.

// RemoveDuplicates removes duplicates in-place and returns the new length.
// Time: O(n)  Space: O(1)
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	return slow + 1
}
