package arrays_slices

// ============================================================
// SOLUTIONS — 08 Arrays & Slices
// ============================================================

// ============================================================
// PART A — Algorithm Pattern Solutions
// ============================================================

func reverse(s []int, lo, hi int) {
	for lo < hi {
		s[lo], s[hi] = s[hi], s[lo]
		lo++
		hi--
	}
}

func ReverseSliceSolution(s []int) {
	reverse(s, 0, len(s)-1)
}

func RemoveDuplicatesSolution(s []int) []int {
	if len(s) == 0 {
		return s
	}
	write := 1
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			s[write] = s[i]
			write++
		}
	}
	return s[:write]
}

func Make2DSolution(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	return matrix
}

func RotateLeftSolution(s []int, k int) {
	n := len(s)
	if n == 0 {
		return
	}
	k = k % n
	reverse(s, 0, k-1)
	reverse(s, k, n-1)
	reverse(s, 0, n-1)
}

func FilterSolution(s []int, fn func(int) bool) []int {
	result := make([]int, 0)
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func MergeSortedSolution(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}

// ============================================================
// PART B — Slice Internals Solutions
// ============================================================

func SafeDeleteSolution(s []int, i int) []int {
	result := make([]int, len(s))
	copy(result, s)
	return append(result[:i], result[i+1:]...)
}

func CopySliceSolution(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

func NilVsEmptySolution() ([]int, []int) {
	var nilSlice []int
	emptySlice := []int{}
	return nilSlice, emptySlice
}

func ExtractWithoutLeakSolution(s []int, from, to int) []int {
	result := make([]int, to-from)
	copy(result, s[from:to])
	return result
}

func ObserveGrowthSolution(n int) []int {
	caps := make([]int, 0, n)
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, i)
		caps = append(caps, cap(s))
	}
	return caps
}

func DetachSliceSolution(s []int) []int {
	return s[0:len(s):len(s)]
}