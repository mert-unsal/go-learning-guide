// Package sort_pkg covers the sort package — critical for coding exams.
// Rename to avoid conflict with built-in sort package name.
package sort_pkg

import (
	"fmt"
	"sort"
)

// ============================================================
// 1. SORTING BASIC TYPES
// ============================================================

func DemonstrateBasicSort() {
	// Sort integers
	nums := []int{5, 2, 8, 1, 9, 3}
	sort.Ints(nums)
	fmt.Println("Sorted ints:", nums) // [1 2 3 5 8 9]

	// Sort strings
	strs := []string{"banana", "apple", "cherry", "date"}
	sort.Strings(strs)
	fmt.Println("Sorted strings:", strs)

	// Sort float64
	floats := []float64{3.14, 1.41, 2.71, 0.57}
	sort.Float64s(floats)
	fmt.Println("Sorted floats:", floats)

	// Check if sorted
	fmt.Println("Is sorted?", sort.IntsAreSorted(nums)) // true

	// Reverse sort
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	fmt.Println("Reverse sorted:", nums) // [9 8 5 3 2 1]
}

// ============================================================
// 2. SORT.SLICE — sort anything with a custom comparator
// ============================================================
// sort.Slice(slice, less func(i, j int) bool)
// less(i, j) returns true if element at i should come BEFORE j

type Person struct {
	Name string
	Age  int
}

func DemonstrateSortSlice() {
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Carol", 35},
		{"Dave", 25},
	}

	// Sort by age ascending
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Println("By age:", people)

	// Sort by name alphabetically
	sort.Slice(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})
	fmt.Println("By name:", people)

	// Sort by age, then by name (stable multi-key sort)
	sort.SliceStable(people, func(i, j int) bool {
		if people[i].Age != people[j].Age {
			return people[i].Age < people[j].Age
		}
		return people[i].Name < people[j].Name
	})
	fmt.Println("By age then name:", people)

	// Sort 2D slice (e.g., intervals by start time)
	intervals := [][]int{{3, 5}, {1, 4}, {2, 6}, {1, 2}}
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i][0] != intervals[j][0] {
			return intervals[i][0] < intervals[j][0]
		}
		return intervals[i][1] < intervals[j][1]
	})
	fmt.Println("Sorted intervals:", intervals)
}

// ============================================================
// 3. SORT.SEARCH — binary search
// ============================================================
// sort.Search(n, f) returns the smallest index i in [0, n)
// at which f(i) is true. f must be monotone (false...false...true...true).

func DemonstrateBinarySearch() {
	nums := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}

	target := 6
	// Find first index where nums[i] >= target
	i := sort.SearchInts(nums, target)
	if i < len(nums) && nums[i] == target {
		fmt.Printf("Found %d at index %d\n", target, i)
	} else {
		fmt.Printf("%d not found\n", target)
	}

	// General sort.Search
	target = 15
	j := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= target
	})
	fmt.Printf("First index where nums[i] >= %d: %d (value: %d)\n", target, j, nums[j])
}

// ============================================================
// 4. IMPLEMENTING sort.Interface (for custom types)
// ============================================================
// To use sort.Sort(x), x must implement:
//   Len() int
//   Less(i, j int) bool
//   Swap(i, j int)

type ByLength []string

func (s ByLength) Len() int           { return len(s) }
func (s ByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByLength) Less(i, j int) bool { return len(s[i]) < len(s[j]) }

func DemonstrateCustomSort() {
	words := []string{"peach", "kiwi", "apple", "fig", "cherry", "banana"}
	sort.Sort(ByLength(words))
	fmt.Println("Sorted by length:", words)

	// Equivalent using sort.Slice (simpler, preferred)
	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) < len(words[j])
	})
	fmt.Println("Sorted by length (Slice):", words)
}

// RunAll runs all demonstrations
func RunAll() {
	fmt.Println("\n=== Basic Sort ===")
	DemonstrateBasicSort()
	fmt.Println("\n=== sort.Slice ===")
	DemonstrateSortSlice()
	fmt.Println("\n=== Binary Search ===")
	DemonstrateBinarySearch()
	fmt.Println("\n=== Custom Sort ===")
	DemonstrateCustomSort()
}
