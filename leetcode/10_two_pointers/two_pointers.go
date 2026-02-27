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

// ============================================================
// PROBLEM 6: Valid Triangle Number (LeetCode #611) — MEDIUM
// ============================================================
// Given an array of non-negative integers, count the number of triplets
// that could form a valid triangle (sum of two sides > third).
//
// Example: nums=[2,2,3,4] → 3
//
// Approach: sort + two pointers.
// Fix the largest side (k), then find pairs (i, j) where nums[i]+nums[j] > nums[k].

// TriangleNumber counts valid triangle triplets.
// Time: O(n²)  Space: O(1)
func TriangleNumber(nums []int) int {
	sort.Ints(nums)
	count := 0
	for k := len(nums) - 1; k >= 2; k-- {
		left, right := 0, k-1
		for left < right {
			if nums[left]+nums[right] > nums[k] {
				count += right - left // all pairs from left to right-1 work
				right--
			} else {
				left++
			}
		}
	}
	return count
}

// ============================================================
// PROBLEM 7: Squares of a Sorted Array (LeetCode #977) — EASY
// ============================================================
// Given a sorted array of integers, return the squares in sorted order.
//
// Example: nums=[-4,-1,0,3,10] → [0,1,9,16,100]
//
// Approach: two pointers from both ends; larger absolute value → larger square.
// Fill result from the back.

// SortedSquares returns sorted squares of a sorted array.
// Time: O(n)  Space: O(n)
func SortedSquares(nums []int) []int {
	n := len(nums)
	result := make([]int, n)
	left, right := 0, n-1
	pos := n - 1

	for left <= right {
		lSq := nums[left] * nums[left]
		rSq := nums[right] * nums[right]
		if lSq > rSq {
			result[pos] = lSq
			left++
		} else {
			result[pos] = rSq
			right--
		}
		pos--
	}
	return result
}

// ============================================================
// PROBLEM 8: 4Sum (LeetCode #18) — MEDIUM
// ============================================================
// Find all unique quadruplets in nums that sum to target.
//
// Example: nums=[1,0,-1,0,-2,2], target=0 → [[-2,-1,1,2],[-2,0,0,2],[-1,0,0,1]]
//
// Approach: sort + two outer loops + two pointers (extension of 3Sum).

// FourSum returns all unique quadruplets summing to target.
// Time: O(n³)  Space: O(1) extra
func FourSum(nums []int, target int) [][]int {
	sort.Ints(nums)
	n := len(nums)
	var result [][]int

	for i := 0; i < n-3; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for j := i + 1; j < n-2; j++ {
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			left, right := j+1, n-1
			for left < right {
				sum := nums[i] + nums[j] + nums[left] + nums[right]
				if sum == target {
					result = append(result, []int{nums[i], nums[j], nums[left], nums[right]})
					for left < right && nums[left] == nums[left+1] {
						left++
					}
					for left < right && nums[right] == nums[right-1] {
						right--
					}
					left++
					right--
				} else if sum < target {
					left++
				} else {
					right--
				}
			}
		}
	}
	return result
}

// ============================================================
// PROBLEM 9: Minimum Difference Between Highest and Lowest of K Scores (LeetCode #1984) — EASY
// ============================================================
// Given an array of scores and k, find the minimum difference between
// the highest and lowest scores among any k students.
//
// Example: nums=[90,100,78,56,70], k=2 → 10  ([90,100])
//
// Approach: sort + sliding window of size k. Min diff = nums[i+k-1] - nums[i].

// MinimumDifference returns the minimum spread of any k elements.
// Time: O(n log n)  Space: O(1)
func MinimumDifference(nums []int, k int) int {
	sort.Ints(nums)
	minDiff := nums[len(nums)-1] - nums[0]
	for i := 0; i+k-1 < len(nums); i++ {
		diff := nums[i+k-1] - nums[i]
		if diff < minDiff {
			minDiff = diff
		}
	}
	return minDiff
}

// ============================================================
// PROBLEM 10: Bag of Tokens (LeetCode #948) — MEDIUM
// ============================================================
// You have tokens[i] power cost. You start with power P and score 0.
// You can play token[i] face-up (spend power, gain score) if P >= tokens[i].
// Or play token[i] face-down (gain 1 power, lose 1 score) if score >= 1.
// Maximize the score.
//
// Example: tokens=[100,200,300,400], power=200 → 2
//
// Approach: sort + two pointers. Greedily face-up cheapest, face-down most expensive.

// BagOfTokensScore returns the maximum achievable score.
// Time: O(n log n)  Space: O(1)
func BagOfTokensScore(tokens []int, power int) int {
	sort.Ints(tokens)
	left, right := 0, len(tokens)-1
	score, maxScore := 0, 0

	for left <= right {
		if power >= tokens[left] {
			power -= tokens[left]
			score++
			left++
			if score > maxScore {
				maxScore = score
			}
		} else if score > 0 {
			power += tokens[right]
			score--
			right--
		} else {
			break
		}
	}
	return maxScore
}
