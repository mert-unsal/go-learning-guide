// For Loops in Go — the ONLY loop keyword.
// Go has ONLY ONE loop keyword: 'for'
// It replaces while, do-while, and for from other languages.
package main

import "fmt"

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  For Loops — Go's ONLY Loop Keyword     %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// Classic C-style for loop
	fmt.Printf("%s▸ Classic C-style: for init; condition; post%s\n", cyan+bold, reset)
	fmt.Print("  ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%s%d%s ", magenta, i, reset)
	}
	fmt.Println()
	fmt.Printf("  %s✔ Identical to C/Java — most familiar form%s\n", green, reset)

	// While-style: just condition (no init or post)
	fmt.Printf("\n%s▸ While-style: for <condition>%s\n", cyan+bold, reset)
	fmt.Printf("  %sNo 'while' keyword in Go — 'for' with just a condition is the equivalent%s\n", dim, reset)
	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Printf("  doubling 1 until ≥ 100 → n = %s%d%s\n", magenta, n, reset)

	// Infinite loop (exit with 'break')
	fmt.Printf("\n%s▸ Infinite Loop: for { }%s\n", cyan+bold, reset)
	fmt.Printf("  %sReplaces 'while(true)' — cleaner, more intentional%s\n", dim, reset)
	count := 0
	for {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Printf("  looped until count ≥ 3 → count = %s%d%s\n", magenta, count, reset)

	// continue: skip to next iteration
	fmt.Printf("\n%s▸ Continue — skip current iteration%s\n", cyan+bold, reset)
	fmt.Print("  odd numbers 0..9: ")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // skip even numbers
		}
		fmt.Printf("%s%d%s ", magenta, i, reset)
	}
	fmt.Println()

	// range over slice
	fmt.Printf("\n%s▸ Range over Slice — yields (index, value)%s\n", cyan+bold, reset)
	nums := []int{10, 20, 30, 40}
	for index, value := range nums {
		fmt.Printf("  nums[%s%d%s] = %s%d%s\n", cyan, index, reset, magenta, value, reset)
	}
	fmt.Printf("  %s✔ Both index and value are COPIES — modifying 'value' won't change the slice%s\n", green, reset)

	// range over string (iterates over RUNES, not bytes)
	fmt.Printf("\n%s▸ Range over String — iterates RUNES, not bytes%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ Index is byte offset, not rune index — matters for multi-byte chars (UTF-8)%s\n", yellow, reset)
	for i, ch := range "Go!" {
		fmt.Printf("  byte-offset=%s%d%s, rune=%s%c%s\n", cyan, i, reset, magenta, ch, reset)
	}

	// range over map (order is RANDOM)
	fmt.Printf("\n%s▸ Range over Map — order is RANDOMIZED%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ Go intentionally randomizes map iteration — never depend on order%s\n", yellow, reset)
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for key, val := range m {
		fmt.Printf("  %s%s%s: %s%d%s\n", cyan, key, reset, magenta, val, reset)
	}

	// Discard index with _
	fmt.Printf("\n%s▸ Blank Identifier _ — discard unwanted index/value%s\n", cyan+bold, reset)
	sum := 0
	for _, v := range []int{1, 2, 3, 4, 5} {
		sum += v
	}
	fmt.Printf("  sum of [1..5] = %s%d%s\n", magenta, sum, reset)

	// Go 1.22+: range over integer
	fmt.Printf("\n%s▸ Range over Integer (Go 1.22+)%s\n", cyan+bold, reset)
	fmt.Printf("  %sSyntax: for i := range N — iterates 0..N-1%s\n", dim, reset)
	fmt.Print("  ")
	for i := range 5 {
		fmt.Printf("%s%d%s ", magenta, i, reset)
	}
	fmt.Println()
	fmt.Printf("  %s✔ Cleaner than 'for i := 0; i < 5; i++' when you just need the count%s\n", green, reset)

	// Labeled break/continue (for nested loops)
	fmt.Printf("\n%s▸ Labeled Break — exit outer loop from inner%s\n", cyan+bold, reset)
	fmt.Printf("  %sbreak outer exits the labeled loop, not just the innermost%s\n", dim, reset)
	fmt.Print("  ")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer // breaks the OUTER loop
			}
			fmt.Printf("(%s%d%s,%s%d%s) ", cyan, i, reset, cyan, j, reset)
		}
	}
	fmt.Println()
	fmt.Printf("  %s✔ Stopped at (1,1) — labeled break exited both loops%s\n", green, reset)
	fmt.Printf("  %s⚠ Use sparingly — labeled breaks can hurt readability in complex code%s\n", yellow, reset)
}
