// Zero Values in Go — every variable is ALWAYS initialized to its zero value.
// This prevents undefined behavior common in C/C++.
//
// Zero value reference table:
//
//	int, float     → 0
//	bool           → false
//	string         → ""
//	pointer        → nil
//	slice          → nil
//	map            → nil
//	channel        → nil
//	function       → nil
//	interface      → nil
//	struct         → all fields set to their zero values
package main

import "fmt"

func main() {
	var i int
	var f float64
	var b bool
	var s string
	fmt.Printf("int: %d, float: %f, bool: %t, string: %q\n", i, f, b, s)
}
