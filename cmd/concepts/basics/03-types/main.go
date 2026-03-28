// Basic Types in Go — demonstrates the type system, conversions, and strings.
//
// Go is STATICALLY TYPED. Every variable has a fixed type.
//
// Integer types:
//
//	int8, int16, int32, int64  (signed)
//	uint8, uint16, uint32, uint64  (unsigned)
//	int, uint  (platform-dependent: 32 or 64 bit)
//	byte = uint8 (alias)
//	rune = int32 (alias, represents a Unicode code point)
//
// Float types:
//
//	float32, float64
//
// Complex types:
//
//	complex64, complex128
//
// String: immutable sequence of bytes (UTF-8 encoded)
// Bool: true or false
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
	fmt.Printf("%s%s  Go's Type System & Conversions          %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// Integer arithmetic
	fmt.Printf("%s▸ Integer arithmetic%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Go integer division truncates toward zero — no implicit float promotion%s\n", green, reset)

	var i int = 10
	var j int = 3
	fmt.Printf("  10 + 3  = %s%d%s\n", magenta, i+j, reset)
	fmt.Printf("  10 - 3  = %s%d%s\n", magenta, i-j, reset)
	fmt.Printf("  10 * 3  = %s%d%s\n", magenta, i*j, reset)
	fmt.Printf("  10 / 3  = %s%d%s  %s(integer division — truncated, not rounded!)%s\n", magenta, i/j, reset, yellow, reset)
	fmt.Printf("  10 %% 3  = %s%d%s\n\n", magenta, i%j, reset)

	// Type conversion (EXPLICIT — no implicit conversion in Go!)
	fmt.Printf("%s▸ Explicit type conversion%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Go has NO implicit type conversion — you must cast explicitly with T(value)%s\n", green, reset)
	fmt.Printf("  %s⚠ This prevents subtle bugs common in C/C++ (e.g., int + float silent promotion)%s\n", yellow, reset)

	var f float64 = float64(i) / float64(j)
	fmt.Printf("  float64(10) / float64(3) = %s%f%s\n\n", magenta, f, reset)

	// String operations
	fmt.Printf("%s▸ String internals — bytes vs characters%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Go strings are immutable byte slices, always UTF-8 encoded%s\n", green, reset)
	fmt.Printf("  %s⚠ len() returns BYTE count, not character (rune) count!%s\n", yellow, reset)

	s := "Hello, 世界" // Go strings support Unicode
	fmt.Printf("  s = %s%q%s\n", magenta, s, reset)
	fmt.Printf("  len(s) = %s%d%s bytes  %s(each CJK char is 3 bytes in UTF-8)%s\n\n", magenta, len(s), reset, dim, reset)

	// To iterate over characters (runes) use range:
	fmt.Printf("%s▸ Rune iteration with range%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ 'range' on a string decodes UTF-8 and yields runes — index may skip bytes%s\n", green, reset)
	for idx, ch := range s {
		if idx < 5 {
			fmt.Printf("  index=%s%d%s  char=%s%c%s  unicode=U+%s%04X%s\n", magenta, idx, reset, magenta, ch, reset, magenta, ch, reset)
		}
	}
	fmt.Println()

	// byte vs rune
	fmt.Printf("%s▸ byte vs rune%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ byte = uint8 (alias) — a single ASCII value%s\n", green, reset)
	fmt.Printf("  %s✔ rune = int32 (alias) — a Unicode code point, can represent any character%s\n", green, reset)

	var myByte byte = 'A' // single quotes = rune/byte literal
	var myRune rune = '世'
	fmt.Printf("  byte 'A' = %s%d%s  (fits in 1 byte)\n", magenta, myByte, reset)
	fmt.Printf("  rune '世' = %s%d%s  (needs int32 to hold U+4E16)\n\n", magenta, myRune, reset)

	// String to/from byte slice
	fmt.Printf("%s▸ String ↔ []byte conversion%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Strings are immutable — to mutate, convert to []byte, modify, convert back%s\n", green, reset)
	fmt.Printf("  %s⚠ Each conversion copies the data — O(n) cost. In hot paths, consider bytes.Buffer%s\n", yellow, reset)

	str := "hello"
	bytes := []byte(str) // convert string to []byte
	bytes[0] = 'H'
	fmt.Printf("  []byte(%q) → modify [0]='H' → %s%s%s\n", str, magenta, string(bytes), reset)
}
