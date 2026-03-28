// Constants in Go — demonstrates const, iota, and bitwise flag patterns.
//
// Constants are declared with 'const' and CANNOT be changed.
// They must be known at compile time.
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

const MaxSize = 100
const AppName = "LearnGo"
const Pi = 3.14159265358979

// Weekday --- iota: auto-incrementing constant generator ---
type Weekday int

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)

// iota with bit shifting — very common for flags/permissions
type Permission uint

const (
	Read    Permission = 1 << iota // 1  (001)
	Write                          // 2  (010)
	Execute                        // 4  (100)
)

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Constants, Iota & Bit Flags             %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Package-level constants ---
	fmt.Printf("%s▸ Package-level constants%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Constants are resolved at compile time — they have no memory address%s\n", green, reset)
	fmt.Printf("  %s✔ Untyped constants have arbitrary precision until assigned to a variable%s\n", green, reset)
	fmt.Printf("  MaxSize = %s%d%s\n", magenta, MaxSize, reset)
	fmt.Printf("  AppName = %s%s%s\n", magenta, AppName, reset)
	fmt.Printf("  Pi      = %s%v%s  %s(untyped float — higher precision than float64 until used)%s\n\n", magenta, Pi, reset, dim, reset)

	// --- iota enumeration ---
	fmt.Printf("%s▸ iota — auto-incrementing enumerator%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ iota resets to 0 in each const block and increments per line%s\n", green, reset)
	fmt.Printf("  %s✔ Subsequent constants repeat the expression with the next iota value%s\n", green, reset)
	fmt.Printf("  Sunday    = %s%d%s  (iota=0)\n", magenta, Sunday, reset)
	fmt.Printf("  Monday    = %s%d%s  (iota=1)\n", magenta, Monday, reset)
	fmt.Printf("  Tuesday   = %s%d%s  (iota=2)\n", magenta, Tuesday, reset)
	fmt.Printf("  Wednesday = %s%d%s  (iota=3)\n", magenta, Wednesday, reset)
	fmt.Printf("  Thursday  = %s%d%s  (iota=4)\n", magenta, Thursday, reset)
	fmt.Printf("  Friday    = %s%d%s  (iota=5)\n", magenta, Friday, reset)
	fmt.Printf("  Saturday  = %s%d%s  (iota=6)\n\n", magenta, Saturday, reset)

	// --- Bit-shifted flags ---
	fmt.Printf("%s▸ iota with bit shifting — permission flags%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ 1 << iota produces powers of 2 — perfect for bitwise flags%s\n", green, reset)
	fmt.Printf("  %s⚠ This is the idiomatic Go pattern for flag sets (like Unix file permissions)%s\n", yellow, reset)
	fmt.Printf("  Read    = 1 << 0 = %s%d%s  (binary: %s%03b%s)\n", magenta, Read, reset, magenta, Read, reset)
	fmt.Printf("  Write   = 1 << 1 = %s%d%s  (binary: %s%03b%s)\n", magenta, Write, reset, magenta, Write, reset)
	fmt.Printf("  Execute = 1 << 2 = %s%d%s  (binary: %s%03b%s)\n\n", magenta, Execute, reset, magenta, Execute, reset)

	// Checking permissions with bitwise AND
	fmt.Printf("%s▸ Bitwise permission checks%s\n", cyan+bold, reset)
	myPerms := Read | Write // 3 (011)
	fmt.Printf("  myPerms = Read | Write = %s%d%s  (binary: %s%03b%s)\n", magenta, myPerms, reset, magenta, myPerms, reset)
	fmt.Printf("  %s✔ Use bitwise AND to test if a flag is set%s\n", green, reset)
	fmt.Printf("  myPerms & Read    != 0 → %s%t%s  (bit 0 is set)\n", magenta, myPerms&Read != 0, reset)
	fmt.Printf("  myPerms & Execute != 0 → %s%t%s  (bit 2 is not set)\n", magenta, myPerms&Execute != 0, reset)
}
