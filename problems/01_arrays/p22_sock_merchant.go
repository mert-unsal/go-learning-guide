package arrays

// ============================================================
// Sales by Match — [E]
// ============================================================
// Given an array of sock colors, find the number of matching pairs.
//
// Example: ar=[10,20,20,10,10,30,50,10,20] → 3 pairs (10,10,20)

// SockMerchant returns the number of matching sock pairs.
// Time: O(n)  Space: O(n)
func SockMerchant(ar []int) int {
	count := make(map[int]int)
	pairs := 0
	for _, sock := range ar {
		count[sock]++
		if count[sock]%2 == 0 {
			pairs++
		}
	}
	return pairs
}
