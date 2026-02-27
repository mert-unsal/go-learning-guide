// Package hackerrank contains problems from HackerRank with full solutions.
// Topics span: arrays, strings, math, sorting, recursion, data structures, greedy.
// Difficulty markers: [E]=Easy, [M]=Medium, [H]=Hard
package hackerrank

import (
	"math"
	"sort"
	"strings"
)

// ============================================================
// PROBLEM 1: Mini-Max Sum (HackerRank) — [E]
// ============================================================
// Given 5 integers, find the minimum and maximum sums of 4 out of 5 numbers.
//
// Example: arr=[1,2,3,4,5]
// minSum = 1+2+3+4 = 10
// maxSum = 2+3+4+5 = 14
// Output: "10 14"
//
// Key insight: sort → minSum = sum of first 4, maxSum = sum of last 4.
// Or: total - max = minSum, total - min = maxSum

// MiniMaxSum returns the minimum and maximum 4-element sums.
// Time: O(n log n)  Space: O(1)
func MiniMaxSum(arr []int) (minSum, maxSum int) {
	sort.Ints(arr)
	for i := 0; i < 4; i++ {
		minSum += arr[i]
	}
	for i := 1; i < 5; i++ {
		maxSum += arr[i]
	}
	return
}

// ============================================================
// PROBLEM 2: FizzBuzz (HackerRank) — [E]
// ============================================================
// For each number from 1 to n:
// - "Fizz" if divisible by 3
// - "Buzz" if divisible by 5
// - "FizzBuzz" if divisible by both
// - the number itself otherwise

// FizzBuzz returns the FizzBuzz sequence up to n.
// Time: O(n)  Space: O(n)
func FizzBuzz(n int) []string {
	result := make([]string, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			result[i-1] = "FizzBuzz"
		case i%3 == 0:
			result[i-1] = "Fizz"
		case i%5 == 0:
			result[i-1] = "Buzz"
		default:
			result[i-1] = intToStr(i)
		}
	}
	return result
}

// ============================================================
// PROBLEM 3: Diagonal Difference (HackerRank) — [E]
// ============================================================
// Given a square matrix, compute the absolute difference
// between the sums of its diagonals.
//
// Example:
//   1 2 3       Primary:   1+5+9 = 15
//   4 5 6       Secondary: 3+5+7 = 15
//   7 8 9       |15 - 15| = 0

// DiagonalDifference returns the absolute diagonal difference.
// Time: O(n)  Space: O(1)
func DiagonalDifference(matrix [][]int) int {
	n := len(matrix)
	primary, secondary := 0, 0
	for i := 0; i < n; i++ {
		primary += matrix[i][i]
		secondary += matrix[i][n-1-i]
	}
	diff := primary - secondary
	if diff < 0 {
		return -diff
	}
	return diff
}

// ============================================================
// PROBLEM 4: Counting Valleys (HackerRank) — [E]
// ============================================================
// A hike is described as a string of U (up) and D (down) steps.
// A valley is a sequence of steps below sea level starting and ending at sea level.
// Count the number of valleys.
//
// Example: "UDDDUDUU" → 1 valley

// CountingValleys counts the number of valleys in the hike.
// Time: O(n)  Space: O(1)
func CountingValleys(steps string) int {
	level := 0
	valleys := 0
	for _, step := range steps {
		if step == 'U' {
			level++
			if level == 0 { // just crossed back to sea level from below
				valleys++
			}
		} else {
			level--
		}
	}
	return valleys
}

// ============================================================
// PROBLEM 5: Sales by Match (HackerRank) — [E]
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

// ============================================================
// PROBLEM 6: Jumping on the Clouds (HackerRank) — [E]
// ============================================================
// Jump on clouds. 0 = safe cloud, 1 = thunder cloud (skip).
// You can jump 1 or 2 clouds. Find the minimum jumps to reach the last cloud.
//
// Example: c=[0,0,1,0,0,1,0] → 4

// JumpingOnClouds returns the minimum number of jumps to the last cloud.
// Time: O(n)  Space: O(1)
func JumpingOnClouds(c []int) int {
	jumps := 0
	i := 0
	n := len(c)
	for i < n-1 {
		// Prefer jump of 2 if it's safe, else jump 1
		if i+2 < n && c[i+2] == 0 {
			i += 2
		} else {
			i++
		}
		jumps++
	}
	return jumps
}

// ============================================================
// PROBLEM 7: Repeated String (HackerRank) — [E]
// ============================================================
// Infinite string formed by repeating s. Count occurrences of 'a' in first n chars.
//
// Example: s="aba", n=10 → 7  ("abaabaabaaba"[:10] = "abaabaabaab" → 7 a's)

// RepeatedString counts 'a' in the first n characters of the infinite repeated string.
// Time: O(|s|)  Space: O(1)
func RepeatedString(s string, n int) int {
	countInS := strings.Count(s, "a")
	fullRepeats := n / len(s)
	remainder := n % len(s)
	countInRemainder := strings.Count(s[:remainder], "a")
	return fullRepeats*countInS + countInRemainder
}

// ============================================================
// PROBLEM 8: Encryption (HackerRank) — [M]
// ============================================================
// Encrypt a string by arranging it in a grid, then reading columns top-to-bottom.
// Grid size: floor(sqrt(len)) rows × ceil(sqrt(len)) cols.
//
// Example: s="haveaniceday" → "hae and via ecy"  (without spaces)
// Grid (3×4):
//   h a v e
//   a n i c
//   e d a y
// Columns: "hae" "and" "via" "ecy" → "haeandviacey"

// Encryption encrypts a string using the grid column method.
// Time: O(n)  Space: O(n)
func Encryption(s string) string {
	// Remove spaces
	s = strings.ReplaceAll(s, " ", "")
	n := len(s)
	rows := int(math.Floor(math.Sqrt(float64(n))))
	cols := int(math.Ceil(math.Sqrt(float64(n))))
	if rows*cols < n {
		rows++
	}

	var result strings.Builder
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			idx := r*cols + c
			if idx < n {
				result.WriteByte(s[idx])
			}
		}
		if c < cols-1 {
			result.WriteByte(' ')
		}
	}
	return result.String()
}

// ============================================================
// PROBLEM 9: Caesar Cipher (HackerRank) — [E]
// ============================================================
// Shift each letter by k positions (wrapping around), preserve case and non-letters.
//
// Example: s="middle-Outz", k=2 → "okffng-Qwvb"

// CaesarCipher applies Caesar cipher with shift k.
// Time: O(n)  Space: O(n)
func CaesarCipher(s string, k int) string {
	k = k % 26 // handle large shifts
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch >= 'a' && ch <= 'z':
			result[i] = byte((int(ch-'a')+k)%26) + 'a'
		case ch >= 'A' && ch <= 'Z':
			result[i] = byte((int(ch-'A')+k)%26) + 'A'
		default:
			result[i] = ch // non-letter: unchanged
		}
	}
	return string(result)
}

// ============================================================
// PROBLEM 10: Pangrams (HackerRank) — [E]
// ============================================================
// A pangram contains every letter of the alphabet at least once.
// Determine if the sentence is a pangram.
//
// Example: "We promptly judged antique ivory buckles for the next prize" → "pangram"

// Pangram returns "pangram" if every letter appears, else "not pangram".
// Time: O(n)  Space: O(1)
func Pangram(sentence string) string {
	var seen [26]bool
	for _, ch := range strings.ToLower(sentence) {
		if ch >= 'a' && ch <= 'z' {
			seen[ch-'a'] = true
		}
	}
	for _, v := range seen {
		if !v {
			return "not pangram"
		}
	}
	return "pangram"
}

// ============================================================
// PROBLEM 11: Array Manipulation (HackerRank) — [H]
// ============================================================
// Given an n-sized array of zeros and m operations, each operation adds value v
// to all elements between indices a and b (1-indexed inclusive).
// Find the maximum value after all operations.
//
// Example: n=5, queries=[[1,2,100],[2,5,100],[3,4,100]] → 200
//
// Key insight: difference array technique.
// Instead of updating all elements in range (O(n) per op),
// add v at index a, subtract v at index b+1 (O(1) per op).
// Then compute prefix sum → original array values.

// ArrayManipulation returns the max value after all range-add operations.
// Time: O(n + m)  Space: O(n)
func ArrayManipulation(n int, queries [][]int) int64 {
	diff := make([]int64, n+2) // 1-indexed, extra slot for n+1

	for _, q := range queries {
		a, b, v := q[0], q[1], int64(q[2])
		diff[a] += v
		if b+1 <= n {
			diff[b+1] -= v
		}
	}

	var maxVal, running int64
	for i := 1; i <= n; i++ {
		running += diff[i]
		if running > maxVal {
			maxVal = running
		}
	}
	return maxVal
}

// ============================================================
// PROBLEM 12: Mark and Toys (HackerRank) — [E]
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

// ============================================================
// PROBLEM 13: Sherlock and the Valid String (HackerRank) — [M]
// ============================================================
// A string is "valid" if all characters have the same frequency.
// OR if removing exactly ONE character makes all frequencies equal.
//
// Example: "aabbcc" → "YES" (all have freq 2)
// Example: "aabbccc" → "YES" (remove one 'c')
// Example: "aabbccdd" → "YES"
// Example: "aabbccd" → "YES"
// Example: "abcde" → "YES" (freq 1 for all)
// Example: "aabbc" → "YES" (remove one 'c')

// IsValid returns "YES" if the string can be made valid by removing at most 1 char.
// Time: O(n)  Space: O(n)
func IsValid(s string) string {
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}
	// Count how many characters have each frequency
	freqOfFreq := make(map[int]int)
	for _, f := range freq {
		freqOfFreq[f]++
	}

	// All same frequency
	if len(freqOfFreq) == 1 {
		return "YES"
	}
	// Two distinct frequencies
	if len(freqOfFreq) == 2 {
		freqs := make([]int, 0, 2)
		counts := make([]int, 0, 2)
		for f, c := range freqOfFreq {
			freqs = append(freqs, f)
			counts = append(counts, c)
		}
		f1, f2 := freqs[0], freqs[1]
		c1, c2 := counts[0], counts[1]
		// One char appears once at freq+1: remove that one char
		if (f1 == f2+1 && c1 == 1) || (f2 == f1+1 && c2 == 1) {
			return "YES"
		}
		// One char appears once at frequency 1 (remove it entirely)
		if (f1 == 1 && c1 == 1) || (f2 == 1 && c2 == 1) {
			return "YES"
		}
		// All chars have same freq and one has freq+1 with only 1 remaining
		if (f1 == f2+1 && c2 == len(freq)-1) || (f2 == f1+1 && c1 == len(freq)-1) {
			return "YES"
		}
	}
	return "NO"
}

// ============================================================
// PROBLEM 14: Climbing the Leaderboard (HackerRank) — [M]
// ============================================================
// Given a leaderboard of scores (descending, with dense ranking) and
// a list of a player's scores, find their rank after each game.
// Dense rank = no gaps (1st, 2nd, 2nd, 3rd, NOT 1st, 2nd, 2nd, 4th).
//
// Example: ranked=[100,100,50,40,40,20,10], player=[5,25,50,120]
// Ranks: [6, 4, 2, 1]
//
// Approach: deduplicate ranked, binary search for each player score.

// ClimbingLeaderboard returns the player's rank after each score.
// Time: O((n + m) log n)  Space: O(n)
func ClimbingLeaderboard(ranked []int, player []int) []int {
	// Deduplicate ranked scores (already sorted descending)
	unique := []int{ranked[0]}
	for i := 1; i < len(ranked); i++ {
		if ranked[i] != ranked[i-1] {
			unique = append(unique, ranked[i])
		}
	}
	n := len(unique)
	result := make([]int, len(player))

	for i, score := range player {
		// Binary search in DESCENDING array: find first index where unique[mid] <= score
		// Rank = position in 1-based index + 1 for scores strictly above
		lo, hi := 0, n-1
		rank := n + 1 // default: after everyone
		for lo <= hi {
			mid := lo + (hi-lo)/2
			if unique[mid] <= score {
				rank = mid + 1 // player ties or beats unique[mid]
				hi = mid - 1   // try to find a better (earlier) position
			} else {
				lo = mid + 1
			}
		}
		result[i] = rank
	}
	return result
}

// ============================================================
// PROBLEM 15: Almost Sorted (HackerRank) — [M]
// ============================================================
// Determine if a permutation can be sorted by:
// - swapping exactly two elements (output "swap a b"), or
// - reversing exactly one contiguous subarray (output "reverse a b")
// Otherwise output "no".
// All indices are 1-based.

// AlmostSorted determines what single operation sorts the array.
func AlmostSorted(arr []int) string {
	n := len(arr)
	// Find the leftmost position that's out of order
	left := -1
	for i := 0; i < n-1; i++ {
		if arr[i] > arr[i+1] {
			left = i
			break
		}
	}
	if left == -1 {
		return "yes" // already sorted
	}
	// Find the rightmost position that's out of order
	right := -1
	for i := n - 1; i > 0; i-- {
		if arr[i] < arr[i-1] {
			right = i
			break
		}
	}
	// Try reversing arr[left..right] and check if fully sorted
	reversed := make([]int, n)
	copy(reversed, arr)
	for i, j := left, right; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	for i := 0; i < n-1; i++ {
		if reversed[i] > reversed[i+1] {
			return "no"
		}
	}
	// Sorted after reversal
	if left+1 == right {
		return "swap " + intToStr(left+1) + " " + intToStr(right+1)
	}
	return "reverse " + intToStr(left+1) + " " + intToStr(right+1)
}

// ============================================================
// HELPER
// ============================================================

// intToStr converts an int to a decimal string (no strconv import needed).
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
