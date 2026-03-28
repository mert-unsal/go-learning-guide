// For Loops in Go — the ONLY loop keyword.
// Go has ONLY ONE loop keyword: 'for'
// It replaces while, do-while, and for from other languages.
package main

import "fmt"

func main() {
	// Classic C-style for loop
	for i := 0; i < 5; i++ {
		fmt.Print(i, " ") // 0 1 2 3 4
	}
	fmt.Println()

	// While-style: just condition (no init or post)
	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Println("n =", n) // 128

	// Infinite loop (exit with 'break')
	count := 0
	for {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Println("count =", count)

	// continue: skip to next iteration
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // skip even numbers
		}
		fmt.Print(i, " ") // 1 3 5 7 9
	}
	fmt.Println()

	// range over slice
	nums := []int{10, 20, 30, 40}
	for index, value := range nums {
		fmt.Printf("nums[%d] = %d\n", index, value)
	}

	// range over string (iterates over RUNES, not bytes)
	for i, ch := range "Go!" {
		fmt.Printf("index=%d, char=%c\n", i, ch)
	}

	// range over map (order is RANDOM)
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for key, val := range m {
		fmt.Printf("%s: %d\n", key, val)
	}

	// Discard index with _
	sum := 0
	for _, v := range []int{1, 2, 3, 4, 5} {
		sum += v
	}
	fmt.Println("sum =", sum)

	// Go 1.22+: range over integer
	for i := range 5 { // 0, 1, 2, 3, 4
		fmt.Print(i, " ")
	}
	fmt.Println()

	// Labeled break/continue (for nested loops)
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer // breaks the OUTER loop
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}
