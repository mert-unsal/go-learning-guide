package escape

// ============================================================
// EXERCISES -- 02 escape analysis: Stack vs Heap Allocation
// ============================================================
// 12 exercises that teach you to think like the Go compiler's
// escape analysis pass. Each function has a specific allocation
// behavior — your job is to implement them AND predict whether
// values escape to heap.
//
// After implementing, verify with:
//   go build -gcflags='-m' ./exercises/advanced/02_escape_analysis/

import (
	"bytes"
	"sync"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: StackOnly -- return a value (no escape)
// ────────────────────────────────────────────────────────────
// Return the sum of two ints. Nothing should escape.

func StackOnly(a, b int) int {
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 2: EscapeToHeap -- return a pointer (forces escape)
// ────────────────────────────────────────────────────────────
// Create an int with value a+b, return a pointer to it.
// The int MUST escape to heap because the pointer outlives the stack frame.

func EscapeToHeap(a, b int) *int {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 3: NoEscapeSlice -- slice stays on stack
// ────────────────────────────────────────────────────────────
// Create a slice of 3 ints, sum them, return the sum.
// If the slice doesn't escape the function, it stays on stack.

func NoEscapeSlice() int {
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 4: EscapeSlice -- slice escapes via return
// ────────────────────────────────────────────────────────────
// Create and return a slice. The backing array escapes to heap.

func EscapeSlice(n int) []int {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 5: InterfaceEscape -- boxing forces escape
// ────────────────────────────────────────────────────────────
// Convert an int to interface{} and return it.
// Values stored in interfaces generally escape because the
// runtime needs a stable pointer for the interface data field.

func InterfaceEscape(n int) interface{} {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 6: ClosureCapture -- closure captures variable
// ────────────────────────────────────────────────────────────
// Return a closure that increments and returns a counter.
// The counter variable escapes because the closure outlives the frame.

func ClosureCapture() func() int {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 7: PreallocateVsAppend -- reduce allocations
// ────────────────────────────────────────────────────────────
// Given n, create a slice of squares [0, 1, 4, 9, ...].
// Use make([]int, n) and index assignment (not append) to avoid
// growslice allocations.

func PreallocateVsAppend(n int) []int {
	return nil
}

// ────────────────────────────────────────────────────────────
// Exercise 8: StringBuilderVsSprintf -- avoid fmt allocations
// ────────────────────────────────────────────────────────────
// Join strings with separator. Use strings.Builder (not fmt.Sprintf).
// strings.Builder avoids the reflect-heavy fmt machinery.

func JoinStrings(parts []string, sep string) string {
	return ""
}

// ────────────────────────────────────────────────────────────
// Exercise 9: StructValue -- pass struct by value (no escape)
// ────────────────────────────────────────────────────────────
// Create a Point{X, Y}, return its distance from origin.
// Pass by value to keep Point on stack.

type Point struct{ X, Y float64 }

func Distance(x, y float64) float64 {
	return 0
}

// ────────────────────────────────────────────────────────────
// Exercise 10: PooledBuffer -- reuse allocations with sync.Pool
// ────────────────────────────────────────────────────────────
// Use a sync.Pool of *bytes.Buffer to format strings without
// allocating a new buffer each time.
// FormatRecord should get a buffer from pool, write to it,
// get the result, then put the buffer back.

var bufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func FormatRecord(name string, value int) string {
	// TODO: get from pool, buf.Reset(), write, result := buf.String(), put back
	return ""
}

// ────────────────────────────────────────────────────────────
// Exercise 11: SliceTricks -- avoid allocation with append tricks
// ────────────────────────────────────────────────────────────
// Remove element at index i from slice WITHOUT allocating a new slice.
// Use the append(s[:i], s[i+1:]...) trick.
// Modifies the original slice. Returns the shortened slice.

func RemoveAt(s []int, i int) []int {
	return s
}

// ────────────────────────────────────────────────────────────
// Exercise 12: ZeroCopyString -- convert []byte to string without copy
// ────────────────────────────────────────────────────────────
// In hot paths, you might want to compare a []byte with a string
// without allocating. Implement a function that checks if a byte
// slice equals a string. Use direct byte comparison (not string(b)).
// Note: string(b) allocates because strings are immutable.

func BytesEqualString(b []byte, s string) bool {
	return false
}
