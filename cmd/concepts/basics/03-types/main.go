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

func main() {
	// Integer arithmetic
	var i int = 10
	var j int = 3
	fmt.Println(i+j, i-j, i*j, i/j, i%j) // 13 7 30 3 1

	// Type conversion (EXPLICIT — no implicit conversion in Go!)
	var f float64 = float64(i) / float64(j)
	fmt.Println(f) // 3.3333...

	// String operations
	s := "Hello, 世界"    // Go strings support Unicode
	fmt.Println(len(s)) // len() returns BYTES, not characters!

	// To iterate over characters (runes) use range:
	for i, ch := range s {
		if i < 5 {
			fmt.Printf("index=%d char=%c unicode=%d\n", i, ch, ch)
		}
	}

	// byte vs rune
	var myByte byte = 'A' // single quotes = rune/byte literal
	var myRune rune = '世'
	fmt.Printf("byte: %d, rune: %d\n", myByte, myRune)

	// String to/from byte slice
	str := "hello"
	bytes := []byte(str) // convert string to []byte
	bytes[0] = 'H'
	fmt.Println(string(bytes)) // "Hello"
}
