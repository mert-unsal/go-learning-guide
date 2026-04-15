package control_flow

// ============================================================
// EXERCISES — 02 Control Flow
// ============================================================
// Implement each function. Run: go test -race -v ./exercises/fundamentals/02_control_flow/
//
// Deep dive: learnings/22_control_flow_under_the_hood.md

// Exercise 1:
// Return "Fizz" if n divisible by 3, "Buzz" if by 5,
// "FizzBuzz" if both, else the number as a string.
// Use switch, NOT if/else.
func FizzBuzzSwitch(n int) string {
	return ""
}

// Exercise 2:
// Sum all integers from 1 to n (inclusive) using a for loop.
func SumTo(n int) int {
	return 0
}

// Exercise 3:
// Return the number of vowels (a,e,i,o,u) in s (case-insensitive).
func CountVowels(s string) int {
	return 0
}

// Exercise 4:
// Return true if n is prime. A prime number is only divisible by 1 and itself.
// Use a for loop with an early return (break/return inside loop).
func IsPrime(n int) bool {
	return false
}

// Exercise 5:
// Use defer to demonstrate execution order.
// Return a slice of strings showing the order that
// "first", "second", "third" would be printed if deferred.
// Hint: defers run LIFO (last in, first out).
func DeferOrder() []string {
	return nil
}

// ─── ADVANCED EXERCISES (Go-Specific) ────────────────────────

// Exercise 6:
// Demonstrate defer + named return interaction.
// Return n * 2, but use a named return and a defer to ADD 10
// to the result before the caller receives it.
//
// Example: DeferModifyReturn(5) → 5*2 + 10 = 20
//
// Under the hood: defer runs after the return value is set but
// before the caller gets it. A closure in defer can modify named returns.
// See: learnings/22, Section 3
func DeferModifyReturn(n int) (result int) {
	return 0
}

// Exercise 7:
// Demonstrate that defer arguments are evaluated at defer time, not execution time.
// Create a counter starting at 0. Defer a function that captures the counter VALUE
// (not a closure). Increment counter 3 times. Return what the deferred function
// would have captured.
//
// Example: DeferArgCapture() → 0 (because defer captured counter when it was 0)
//
// Under the hood: the compiler evaluates and copies the argument into the _defer
// record at the defer statement. The runtime replays the saved value.
func DeferArgCapture() int {
	return -1
}

// Exercise 8:
// Given a slice of structs, use range to double each player's score IN PLACE.
// The catch: range gives you a COPY of each element, not a reference.
// You must use the index to modify the original.
//
// Under the hood: the compiler rewrites range to v := slice[i], where v
// is a stack-local copy. Modifying v does not affect the slice.
type Player struct {
	Name  string
	Score int
}

func DoubleScores(players []Player) {
}

// Exercise 9:
// Given a string, return a slice of its rune values using range.
// Range over a string iterates RUNES, not bytes.
//
// Example: RuneValues("Go🚀") → []rune{'G', 'o', '🚀'}
//
// Under the hood: the compiler rewrites range-over-string to use
// utf8.DecodeRuneInString, advancing by rune width (1-4 bytes).
func RuneValues(s string) []rune {
	return nil
}

// Exercise 10:
// Search a 2D matrix for a target value. Return its (row, col) position.
// If not found, return (-1, -1).
// You MUST use a labeled break to exit both loops when found.
//
// Under the hood: labeled break exits the named outer loop, not just
// the inner loop. This is idiomatic Go for nested search.
func FindInMatrix(matrix [][]int, target int) (int, int) {
	return -1, -1
}

// Exercise 11:
// Implement a type switch that describes what's inside an interface value.
// Return:
//   int    → "int: <value>"
//   string → "string: <value>"
//   bool   → "bool: <value>"
//   []int  → "slice: len=<len>"
//   nil    → "nil"
//   other  → "unknown"
//
// Under the hood: type switch compares the _type pointer from the
// interface header against known type descriptors. Pointer equality, very fast.
func TypeDescribe(v interface{}) string {
	return ""
}

// Exercise 12:
// Use range over an integer (Go 1.22+) to return a slice of squares
// from 0 to n-1.
//
// Example: Squares(4) → [0, 1, 4, 9]
//
// Under the hood: range N rewrites to for i := 0; i < N; i++.
func Squares(n int) []int {
	return nil
}
