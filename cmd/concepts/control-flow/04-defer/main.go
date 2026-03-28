// Defer in Go — delays execution until the surrounding function returns.
// Very useful for cleanup (closing files, unlocking mutexes, etc.)
//
// KEY RULES:
//  1. Deferred calls run in LIFO order (last deferred = first to run)
//  2. Arguments to deferred functions are evaluated IMMEDIATELY (at defer time)
//  3. Deferred functions CAN read and modify named return values
//  4. Do NOT defer inside loops for resource cleanup — they stack up
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

// readFile — typical defer usage: resource cleanup
func readFile(filename string) {
	// In real code:
	// f, err := os.Open(filename)
	// if err != nil { ... }
	// defer f.Close() ← guaranteed to run when function exits

	fmt.Printf("  1. Opening %s%s%s\n", magenta, filename, reset)
	defer fmt.Printf("  3. Closing %s%s%s ← %s✔ defer ran AFTER function body finished%s\n", magenta, filename, reset, green, reset)
	fmt.Printf("  2. Reading %s%s%s\n", magenta, filename, reset)
}

// deferInLoop — GOTCHA: defer inside a loop — do NOT do this for resources!
// All defers stack up and run only when the FUNCTION exits, not each iteration.
func deferInLoop() {
	fmt.Printf("\n%s▸ Defer in Loop (GOTCHA — resource leak!)%s\n", cyan+bold, reset)
	fmt.Printf("  %s⚠ All defers stack up and run only when the FUNCTION exits, not each iteration%s\n", yellow, reset)
	fmt.Printf("  %s⚠ BAD for file handles / DB connections — they stay open until function returns!%s\n", yellow, reset)
	fmt.Printf("  %sLoop: defer fmt.Println(i) for i = 0, 1, 2%s\n", dim, reset)
	for i := 0; i < 3; i++ {
		// BAD for file handles / DB connections — they stay open until function returns!
		defer fmt.Printf("  defer in loop: i=%s%d%s ← %sLIFO: last deferred runs first%s\n", magenta, i, reset, green, reset)
	}
	fmt.Printf("  %s(deferred prints will appear AFTER main() exits — LIFO order: 2, 1, 0)%s\n", dim, reset)
	// FIX: use an anonymous function or a helper function instead:
	// for i := 0; i < 3; i++ {
	//     func(i int) {
	//         // open resource, defer close, use resource — all scoped here
	//     }(i)
	// }
	fmt.Printf("  %s✔ FIX: wrap loop body in a closure so defer runs each iteration%s\n", green, reset)
}

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Defer in Go — LIFO Cleanup Mechanism   %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	fmt.Printf("%s▸ LIFO Order Demo (defer stack)%s\n", cyan+bold, reset)
	fmt.Printf("  %sDefers are pushed onto a stack — last in, first out%s\n", dim, reset)
	fmt.Printf("  → %sfmt.Println(\"start\")%s\n", dim, reset)
	fmt.Println("  start")
	defer fmt.Printf("  deferred 1 ← %sran LAST (pushed first onto defer stack)%s\n", yellow, reset)
	defer fmt.Printf("  deferred 2 ← %sran second%s\n", yellow, reset)
	defer fmt.Printf("  deferred 3 ← %sran FIRST (pushed last onto defer stack)%s\n", green, reset)
	fmt.Printf("  → %sthree defers registered: 1, 2, 3%s\n", dim, reset)
	fmt.Printf("  → %sfmt.Println(\"end\")%s\n", dim, reset)
	fmt.Println("  end")
	fmt.Printf("  %s✔ Now main() continues… defers will fire in LIFO order at function exit%s\n", green, reset)

	fmt.Printf("\n%s▸ Resource Cleanup Pattern%s\n", cyan+bold, reset)
	fmt.Printf("  %sPattern: open → defer close → use (close is guaranteed even on panic)%s\n", dim, reset)
	readFile("report.txt")

	deferInLoop()

	fmt.Printf("\n%s── main() returning now — all remaining defers fire below ──%s\n", blue+bold, reset)
}
