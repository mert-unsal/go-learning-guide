package basics

// ============================================================
// SOLUTIONS — 01 Basics
// ============================================================

// Solution 1: Celsius to Fahrenheit
// The formula is straightforward. Use float64 arithmetic.
func CelsiusToFahrenheitSolution(c float64) float64 {
	return c*9/5 + 32
}

// Solution 2: Swap two integers
// Go's multiple return values make this clean and idiomatic.
// No need for a temporary variable!
func SwapIntsSolution(a, b int) (int, int) {
	return b, a
	// Or you could do: a, b = b, a; return a, b
}

// Solution 3: Character count (Unicode-aware)
// len(s) counts BYTES. To count Unicode characters, convert to []rune.
// A rune is an int32 representing a Unicode code point.
func CharacterCountSolution(s string) int {
	return len([]rune(s))
	// Example: len("hello") = 5 bytes, len([]rune("hello")) = 5 chars
	// Example: len("世界") = 6 bytes, len([]rune("世界")) = 2 chars
}

// Solution 4: Min and Max from three integers
// Named return values are declared in the signature: (min, max int)
// You can return them without listing them (naked return), but explicit is clearer.
func MinMaxSolution(a, b, c int) (min, max int) {
	// Initialize min/max with the first value
	min, max = a, a

	if b < min {
		min = b
	}
	if b > max {
		max = b
	}
	if c < min {
		min = c
	}
	if c > max {
		max = c
	}
	return min, max
}

// Solution 5: Direction name
// Switch in Go does NOT fall through by default (unlike C/Java).
// You can use 'fallthrough' keyword explicitly if needed.
func DirectionNameSolution(d Direction) string {
	switch d {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return "Unknown"
	}
}
