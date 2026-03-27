package arrays

import "sort"

// ============================================================
// Mark and Toys — [E]
// ============================================================
// Given toy prices and a budget, find the maximum number of toys you can buy.
// You must buy distinct toys (each toy once). Minimize cost per toy.
//
// Example: prices=[1,12,5,111,200,1000,10], k=50 → 4 (buy 1,5,10,12 = 28 ≤ 50)
//
// Greedy: sort by price, buy cheapest first until budget is exhausted.

// MaximumToys returns the maximum number of toys buyable within budget k.
// Time: O(n log n)  Space: O(1)
func MaximumToys(prices []int, k int) int {
	sort.Ints(prices)
	count := 0
	for _, price := range prices {
		if k < price {
			break
		}
		k -= price
		count++
	}
	return count
}
