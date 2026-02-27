// Package builtins covers Go's built-in functions: make, new, copy, append,
// len, cap, delete, close, panic, recover, and type conversions.
//
// These are not imported from any package — they are always available.
package builtins

import "fmt"

// ============================================================
// 1. make() — Allocate and Initialize Slices, Maps, Channels
// ============================================================
// make() is ONLY for slices, maps, and channels.
// It allocates memory AND initializes the internal structure.
//
// Signatures:
//   make([]T, length)              — slice with len=cap=length
//   make([]T, length, capacity)    — slice with len=length, cap=capacity
//   make(map[K]V)                  — empty map
//   make(map[K]V, hint)            — pre-allocated map (hint=expected size)
//   make(chan T)                    — unbuffered channel
//   make(chan T, capacity)          — buffered channel

func DemonstrateMake() {
	// --- Slice with make ---
	s1 := make([]int, 5)     // len=5, cap=5, all zeros: [0 0 0 0 0]
	s2 := make([]int, 3, 10) // len=3, cap=10: [0 0 0]  (room for 10)
	fmt.Println("s1:", s1, "len:", len(s1), "cap:", cap(s1))
	fmt.Println("s2:", s2, "len:", len(s2), "cap:", cap(s2))

	// Why pre-allocate? Avoids repeated reallocations when appending.
	// If you know you'll append ~1000 items, start with make([]int, 0, 1000).
	efficient := make([]int, 0, 1000) // len=0, cap=1000 — no reallocation until 1001
	for i := 0; i < 5; i++ {
		efficient = append(efficient, i*i)
	}
	fmt.Println("efficient:", efficient) // [0 1 4 9 16]

	// --- Map with make ---
	m := make(map[string]int) // empty map, ready to use
	m["alice"] = 30
	m["bob"] = 25
	fmt.Println("map:", m)

	m2 := make(map[string]int, 100) // hint: expect ~100 entries (less rehashing)
	_ = m2

	// --- Channel with make ---
	ch := make(chan int, 5) // buffered: send up to 5 without a receiver
	ch <- 42
	val := <-ch
	fmt.Println("channel val:", val)
}

// ============================================================
// 2. new() — Allocate Zeroed Memory, Return Pointer
// ============================================================
// new(T) allocates memory for type T, zeros it, returns *T.
// Less common than make() or composite literals.
// Use new() when you need a pointer to a zero value.
//
// Equivalent: p := new(int)  ↔  var x int; p := &x

func DemonstrateNew() {
	p := new(int) // *int pointing to 0
	*p = 42
	fmt.Println("new int:", *p) // 42

	type Point struct{ X, Y int }
	pt := new(Point) // *Point pointing to {0, 0}
	pt.X = 10
	pt.Y = 20
	fmt.Println("new Point:", *pt) // {10 20}

	// In practice, composite literals are more idiomatic:
	pt2 := &Point{X: 10, Y: 20} // same result, more readable
	fmt.Println("literal:", *pt2)
}

// ============================================================
// 3. append() — Add Elements to a Slice
// ============================================================
// append(slice, elem...)  — returns a new (or same) slice
// IMPORTANT: always reassign! append may return a new underlying array.
//
// append grows the capacity automatically (usually doubles when full).

func DemonstrateAppend() {
	var s []int                                  // nil slice — safe to append to
	fmt.Println("nil slice:", s, len(s), cap(s)) // [] 0 0

	// Append single elements
	s = append(s, 1)
	s = append(s, 2, 3, 4)          // append multiple at once
	fmt.Println("after append:", s) // [1 2 3 4]

	// Append another slice with ... (spread operator)
	other := []int{5, 6, 7}
	s = append(s, other...)
	fmt.Println("after spread:", s) // [1 2 3 4 5 6 7]

	// Append to a sub-slice (CAREFUL: may overwrite the original!)
	a := []int{1, 2, 3, 4, 5}
	b := a[:3]                          // b shares underlying array with a
	b = append(b, 99)                   // OVERWRITES a[3]!
	fmt.Println("a after b append:", a) // [1 2 3 99 5]  ← a[3] changed!

	// Safe approach: use copy to avoid shared backing array
	c := make([]int, len(a[:3]))
	copy(c, a[:3])
	c = append(c, 99)                        // safe — c has its own array
	fmt.Println("a after safe c append:", a) // unchanged
}

// ============================================================
// 4. copy() — Copy Between Slices
// ============================================================
// copy(dst, src) copies min(len(dst), len(src)) elements.
// Returns the number of elements copied.
// dst and src may overlap (copy handles this correctly).

func DemonstrateCopy() {
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))

	n := copy(dst, src)
	fmt.Println("copied:", n, "result:", dst) // copied: 5 result: [1 2 3 4 5]

	// Partial copy — copies min(len(dst), len(src))
	small := make([]int, 3)
	copy(small, src)                    // only copies 3 elements
	fmt.Println("partial copy:", small) // [1 2 3]

	// Copy string to []byte
	b := make([]byte, 5)
	copy(b, "Hello")
	fmt.Println("string→[]byte:", b) // [72 101 108 108 111]

	// Shift slice elements left by 1 (delete index 0)
	nums := []int{10, 20, 30, 40, 50}
	copy(nums, nums[1:]) // shift left
	nums = nums[:len(nums)-1]
	fmt.Println("shifted:", nums) // [20 30 40 50]

	// Deep copy a 2D slice — must copy each row individually!
	matrix := [][]int{{1, 2}, {3, 4}}
	matCopy := make([][]int, len(matrix))
	for i, row := range matrix {
		matCopy[i] = make([]int, len(row))
		copy(matCopy[i], row)
	}
	matrix[0][0] = 999
	fmt.Println("original changed:", matrix[0][0])  // 999
	fmt.Println("copy unchanged:  ", matCopy[0][0]) // 1
}

// ============================================================
// 5. len() and cap()
// ============================================================
// len(x) — number of elements currently in x
// cap(x) — total capacity of the underlying array
//
// Works on: strings, arrays, slices, maps, channels

func DemonstrateLenCap() {
	s := make([]int, 3, 8)
	fmt.Printf("len=%d  cap=%d\n", len(s), cap(s)) // len=3 cap=8

	// Watch cap grow as we append beyond current capacity
	s2 := []int{}
	for i := 0; i < 10; i++ {
		s2 = append(s2, i)
		fmt.Printf("len=%d cap=%d\n", len(s2), cap(s2))
		// Cap doubles: 0 → 1 → 2 → 4 → 8 → 16 ...
	}
}

// ============================================================
// 6. delete() — Remove a Key from a Map
// ============================================================
// delete(map, key) — no-op if key doesn't exist (safe to call always)

func DemonstrateDelete() {
	freq := map[string]int{"a": 3, "b": 1, "c": 5}
	fmt.Println("before:", freq)

	delete(freq, "b")           // remove key "b"
	delete(freq, "z")           // no-op — "z" doesn't exist
	fmt.Println("after:", freq) // map[a:3 c:5]

	// Pattern: decrement count, then remove if zero
	delete(freq, "a")
	if freq["a"] == 0 { // zero value returned for missing key
		fmt.Println("'a' is gone")
	}
}

// ============================================================
// 7. Type Conversions — int, string, byte, rune
// ============================================================
// Go requires EXPLICIT conversions — no implicit casting.

func DemonstrateConversions() {
	// int ↔ float64
	var i int = 42
	var f float64 = float64(i)
	var i2 int = int(f) // truncates (does NOT round)
	fmt.Println(i, f, i2)

	// int ↔ string (using rune/byte)
	ch := 'A'         // rune (int32) value: 65
	str := string(ch) // "A"  (interprets as Unicode code point)
	fmt.Println("rune to string:", str)

	// string ↔ []byte
	s := "Hello, 世界"
	bytes := []byte(s)    // encodes to UTF-8 bytes
	back := string(bytes) // decodes from UTF-8
	fmt.Println("bytes:", bytes[:5], "back:", back)

	// string ↔ []rune (Unicode-safe character access)
	runes := []rune(s)
	fmt.Println("rune count:", len(runes))       // 9 (not 13 bytes)
	fmt.Println("third rune:", string(runes[2])) // "l"

	// int → string via strconv (NOT string(int) — that gives a Unicode char!)
	// string(65) = "A", NOT "65"
	// Use: strconv.Itoa(65) = "65"
	fmt.Println("string(65):", string(rune(65))) // "A" — character, not number!

	// Byte arithmetic (very common in LeetCode)
	var b byte = 'z'
	fmt.Println("'z' - 'a' =", b-'a')                // 25
	fmt.Println("'A' + 32  =", string(rune('A'+32))) // "a"

	// int64 ↔ int (be careful on 32-bit platforms)
	var big int64 = 1 << 40
	small := int(big) // safe on 64-bit OS
	fmt.Println("int64 to int:", small)
}

// ============================================================
// 8. panic() and recover() — Error Handling for Exceptional Cases
// ============================================================
// panic: stops normal execution, unwinds the stack (deferred functions still run)
// recover: catches a panic inside a deferred function
//
// Use sparingly! Prefer returning errors. Use panic only for truly unrecoverable states.

// SafeDivide recovers from a divide-by-zero panic and returns an error string.
func SafeDivide(a, b int) (result int, err string) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Sprintf("recovered panic: %v", r)
		}
	}()
	result = a / b // panics if b == 0
	return result, ""
}

func DemonstratePanicRecover() {
	res, err := SafeDivide(10, 2)
	fmt.Println("10/2 =", res, err) // 10/2 = 5

	res2, err2 := SafeDivide(10, 0)
	fmt.Println("10/0 =", res2, err2) // 10/0 = 0 recovered panic: runtime error: integer divide by zero
}
