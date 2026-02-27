// Package arrays contains LeetCode array problems with detailed explanations.
// Topics: hash maps, prefix products, greedy single-pass algorithms.
package arrays

// ============================================================
// PROBLEM 1: Two Sum (LeetCode #1) — EASY
// ============================================================
// Given an array of integers and a target, return indices of the two numbers
// that add up to target. Assume exactly one solution exists.
//
// Example: nums=[2,7,11,15], target=9 → [0,1] (2+7=9)
//
// Brute force: O(n²) — check every pair
// Optimal: O(n) — hash map stores "what complement have we seen?"

// TwoSum returns the indices of two numbers that sum to target.
// Time: O(n)  Space: O(n)
func TwoSum(nums []int, target int) []int {
	// seen maps value → index
	seen := make(map[int]int)

	for i, num := range nums {
		complement := target - num
		if j, ok := seen[complement]; ok {
			return []int{j, i} // complement was seen at index j
		}
		seen[num] = i // record this number's index
	}
	return nil // no solution (problem guarantees one exists)
}

// ============================================================
// PROBLEM 2: Best Time to Buy and Sell Stock (LeetCode #121) — EASY
// ============================================================
// Given daily prices, find the maximum profit from one buy + one sell.
// You must buy BEFORE you sell. Return 0 if no profit is possible.
//
// Example: prices=[7,1,5,3,6,4] → 5 (buy at 1, sell at 6)
//
// Key insight: track the minimum price seen so far.
// At each day, profit = current price - minimum so far.

// MaxProfit returns the maximum profit achievable from one transaction.
// Time: O(n)  Space: O(1)
func MaxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	minPrice := prices[0] // lowest buy price seen so far
	maxProfit := 0

	for _, price := range prices {
		if price < minPrice {
			minPrice = price // found a cheaper buy day
		} else if price-minPrice > maxProfit {
			maxProfit = price - minPrice // better profit today
		}
	}
	return maxProfit
}

// ============================================================
// PROBLEM 3: Product of Array Except Self (LeetCode #238) — MEDIUM
// ============================================================
// Given an array, return an array where output[i] = product of all elements
// EXCEPT nums[i]. Solve in O(n) WITHOUT using division.
//
// Example: nums=[1,2,3,4] → [24,12,8,6]
//
// Key insight: output[i] = (product of everything to the LEFT of i)
//                        × (product of everything to the RIGHT of i)
// Pass 1: fill result with prefix products (left of i).
// Pass 2: sweep right-to-left, multiply by running suffix product.

// ProductExceptSelf returns the product array without division.
// Time: O(n)  Space: O(1) extra (output array doesn't count)
func ProductExceptSelf(nums []int) []int {
	n := len(nums)
	result := make([]int, n)

	// Pass 1: result[i] = product of nums[0..i-1]
	result[0] = 1
	for i := 1; i < n; i++ {
		result[i] = result[i-1] * nums[i-1]
	}

	// Pass 2: multiply result[i] by product of nums[i+1..n-1]
	suffix := 1
	for i := n - 1; i >= 0; i-- {
		result[i] *= suffix
		suffix *= nums[i]
	}
	return result
}

// ============================================================
// PROBLEM 4: Contains Duplicate (LeetCode #217) — EASY
// ============================================================
// Return true if any value appears at least twice.
//
// Time: O(n)  Space: O(n)

// ContainsDuplicate returns true if the slice has any repeated element.
func ContainsDuplicate(nums []int) bool {
	seen := make(map[int]bool)
	for _, n := range nums {
		if seen[n] {
			return true
		}
		seen[n] = true
	}
	return false
}

// ============================================================
// PROBLEM 5: Maximum Subarray (LeetCode #53) — MEDIUM — Kadane's Algorithm
// ============================================================
// Find the contiguous subarray with the largest sum.
//
// Example: nums=[-2,1,-3,4,-1,2,1,-5,4] → 6 ([4,-1,2,1])
//
// Kadane's: at each position, decide: start fresh OR extend previous subarray.
// current = max(nums[i], current + nums[i])

// MaxSubArray returns the largest sum of any contiguous subarray.
// Time: O(n)  Space: O(1)
func MaxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	current := nums[0]
	best := nums[0]

	for _, num := range nums[1:] {
		// Extend previous subarray or start a new one here
		if current+num > num {
			current = current + num
		} else {
			current = num
		}
		if current > best {
			best = current
		}
	}
	return best
}

// ============================================================
// PROBLEM 6: Merge Sorted Array (LeetCode #88) — EASY
// ============================================================
// Merge nums2 into nums1 in-place. nums1 has extra space at the end.
// m = valid elements in nums1, n = valid elements in nums2.
//
// Example: nums1=[1,2,3,0,0,0], m=3, nums2=[2,5,6], n=3 → [1,2,2,3,5,6]
//
// Key insight: fill from the BACK to avoid overwriting elements we still need.
// Compare from the largest ends and place the bigger element at the back.

// Merge merges nums2 into nums1 in-place (fills from the back).
// Time: O(m+n)  Space: O(1)
func Merge(nums1 []int, m int, nums2 []int, n int) {
	i, j, k := m-1, n-1, m+n-1 // i=last of nums1, j=last of nums2, k=write position
	for j >= 0 {
		if i >= 0 && nums1[i] > nums2[j] {
			nums1[k] = nums1[i]
			i--
		} else {
			nums1[k] = nums2[j]
			j--
		}
		k--
	}
}

// ============================================================
// PROBLEM 7: Find All Numbers Disappeared in an Array (LeetCode #448) — EASY
// ============================================================
// Given n integers where each is in [1, n], return all numbers in [1, n]
// that do NOT appear in the array.
//
// Example: nums=[4,3,2,7,8,2,3,1] → [5,6]
//
// Key insight: use the array itself as a hash map.
// For each value v, mark index (|v|-1) negative.
// Then indices still positive correspond to missing numbers.

// FindDisappearedNumbers returns missing numbers using O(1) extra space.
// Time: O(n)  Space: O(1) extra
func FindDisappearedNumbers(nums []int) []int {
	// Mark visited indices negative
	for _, v := range nums {
		if v < 0 {
			v = -v
		}
		idx := v - 1
		if nums[idx] > 0 {
			nums[idx] = -nums[idx]
		}
	}
	// Indices still positive are missing numbers
	var result []int
	for i, v := range nums {
		if v > 0 {
			result = append(result, i+1)
		}
	}
	return result
}

// ============================================================
// PROBLEM 8: Rotate Array (LeetCode #189) — MEDIUM
// ============================================================
// Rotate array to the right by k steps in-place.
//
// Example: nums=[1,2,3,4,5,6,7], k=3 → [5,6,7,1,2,3,4]
//
// Key insight: reverse the whole array, then reverse first k, then rest.
// Reverse all:    [7,6,5,4,3,2,1]
// Reverse [0..k): [5,6,7,4,3,2,1]
// Reverse [k..n): [5,6,7,1,2,3,4]  ✓

// Rotate rotates the array to the right by k positions in-place.
// Time: O(n)  Space: O(1)
func Rotate(nums []int, k int) {
	n := len(nums)
	k = k % n // handle k >= n
	reverse(nums, 0, n-1)
	reverse(nums, 0, k-1)
	reverse(nums, k, n-1)
}

func reverse(nums []int, left, right int) {
	for left < right {
		nums[left], nums[right] = nums[right], nums[left]
		left++
		right--
	}
}

// ============================================================
// PROBLEM 9: Find Minimum in Rotated Sorted Array (LeetCode #153) — MEDIUM
// ============================================================
// (Also in binary_search — here for completeness in arrays track)
// Array was sorted then rotated. Find the minimum element.
//
// Example: nums=[3,4,5,1,2] → 1

// FindMinRotated returns the minimum of a rotated sorted array.
// Time: O(log n)  Space: O(1)
func FindMinRotated(nums []int) int {
	left, right := 0, len(nums)-1
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] > nums[right] {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return nums[left]
}

// ============================================================
// PROBLEM 10: Subarray Sum Equals K (LeetCode #560) — MEDIUM
// ============================================================
// Given an array of integers and k, return the total number of continuous
// subarrays whose sum equals k.
//
// Example: nums=[1,1,1], k=2 → 2
//
// Key insight: prefix sums. If prefixSum[j] - prefixSum[i] == k,
// then subarray [i+1..j] has sum k.
// Use a map: count how many times each prefix sum has been seen.

// SubarraySum returns the number of subarrays with sum equal to k.
// Time: O(n)  Space: O(n)
func SubarraySum(nums []int, k int) int {
	// prefixCount[s] = how many times prefix sum s has occurred
	prefixCount := map[int]int{0: 1} // empty prefix has sum 0
	sum := 0
	count := 0

	for _, num := range nums {
		sum += num
		// If (sum - k) was seen before, those subarrays end here with sum k
		count += prefixCount[sum-k]
		prefixCount[sum]++
	}
	return count
}
