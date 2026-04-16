# 03 — Strings

> How Go represents, stores, and manipulates strings under the hood —
> and why strings are "read-only slices of bytes."
> How understanding string internals directly shapes production logging,
> error handling, and high-throughput system design.

---

## Table of Contents

1. [The String Header — A 2-Word Struct](#1-the-string-header--a-2-word-struct)
2. [Strings Are Immutable — Why and How](#2-strings-are-immutable--why-and-how)
3. [Bytes vs Runes vs Characters](#3-bytes-vs-runes-vs-characters)
4. [String Iteration — Two Different Loops](#4-string-iteration--two-different-loops)
5. [String ↔ Byte Slice Conversions — The Hidden Cost](#5-string--byte-slice-conversions--the-hidden-cost)
6. [Substrings Share Backing Data — Memory Leak Trap](#6-substrings-share-backing-data--memory-leak-trap)
7. [String Concatenation — Performance Analysis](#7-string-concatenation--performance-analysis)
8. [String Comparison and Interning](#8-string-comparison-and-interning)
9. [Strings and the Compiler — Escape Analysis](#9-strings-and-the-compiler--escape-analysis)
   - 9a. [From Strings to Production Logging — A Case Study](#9a-from-strings-to-production-logging--a-case-study)
10. [Production Patterns and Best Practices](#10-production-patterns-and-best-practices)
11. [Quick Reference Card](#11-quick-reference-card)

---

## 1. The String Header — A 2-Word Struct

A string in Go is not an object, not a class — it's a **2-word struct**:

```go
// Source: runtime/string.go
type stringStruct struct {
    str unsafe.Pointer  // pointer to backing byte array (read-only)
    len int             // number of bytes (NOT characters)
}
```

```
string header (16 bytes on 64-bit):
┌──────────────────────┬──────────────────────┐
│  str unsafe.Pointer  │  len int             │
│  (points to bytes)   │  (byte count)        │
└──────────────────────┴──────────────────────┘
```

**Compare with a slice header (24 bytes):**

```
slice header:   { ptr, len, cap }   ← 3 words (24 bytes)
string header:  { ptr, len }        ← 2 words (16 bytes) — no cap, because immutable
```

There's no `cap` field because strings can't grow — they're immutable. There's nothing
to "grow into," so capacity is meaningless.

### What `len()` Returns

```go
s := "Hello"
len(s)  // 5 — five BYTES, not five characters

s = "世界"
len(s)  // 6 — two characters, but 6 BYTES (3 bytes per Chinese character in UTF-8)
```

`len()` reads the `len` field of the string header directly — it's an **O(1)** operation,
inlined by the compiler. It does not count characters.

---

## 2. Strings Are Immutable — Why and How

Once created, a string's backing bytes **cannot be modified**. There is no `s[0] = 'X'`.

```go
s := "hello"
s[0] = 'H'    // ❌ compile error: cannot assign to s[0] (value of type byte)
```

### Why Immutability?

**1. Safe concurrent access.** Multiple goroutines can read the same string without locks.
No data race is possible because nobody can write.

**2. Safe to share backing data.** Substrings and string copies can point to the same
bytes without risk of mutation:

```go
s := "hello world"
sub := s[0:5]     // "hello" — shares the same backing bytes as s
// If strings were mutable, modifying s would corrupt sub
```

**3. Strings can be map keys.** Map keys must be comparable and stable. If a string
could be modified after insertion, the hash would be wrong and lookup would break.

**4. The compiler can optimize.** String literals are stored in the read-only data
segment of the binary. The OS can memory-map them and share across processes.

### How Immutability Is Enforced

The compiler simply **refuses to compile** assignments to string indices. At the runtime
level, the backing bytes for string literals live in the `.rodata` (read-only data)
section of the binary. Attempting to write there via `unsafe` would cause a **segfault**
from the OS memory protection.

### String Pass-By-Value and Goroutine Safety

Go passes the string **header** (16 bytes) by value. But the pointer inside still
points to the **same backing bytes**. This is safe ONLY because strings are immutable:

```
  func process(s string) { ... }

  go process(greeting)

  Main goroutine stack:       New goroutine stack:
  ┌─────────────────┐        ┌─────────────────┐
  │ greeting.ptr ────┼───┐   │ s.ptr ───────────┼───┐
  │ greeting.len = 5 │   │   │ s.len = 5        │   │
  └─────────────────┘   │   └─────────────────┘   │
                         │                          │
                         ▼                          ▼
                  ┌──────────────────────────┐
                  │  "hello" (read-only)     │  ← SHARED backing bytes
                  │  in .rodata segment      │     but nobody can write
                  └──────────────────────────┘
```

Both goroutines hold **separate header copies** pointing to the **same bytes**.
This is a data race in languages with mutable strings — but safe in Go because:

1. **Compiler blocks** `s[0] = 'X'` — won't compile
2. **OS blocks** writes to `.rodata` — hardware-level protection (segfault)
3. **Reassignment** `s = "world"` creates a new header + new backing bytes,
   leaving the other goroutine's view completely untouched

```go
  greeting := "hello"
  go func() {
      greeting = "world"  // ← new header + new backing array allocated
                          //   does NOT modify the "hello" bytes
                          //   but WARNING: this IS a data race on the
                          //   header variable itself (not the bytes)
  }()
```

**Critical distinction:** the backing bytes are always safe to share.
But the header variable (the 16-byte struct on the stack) is still subject to
data races if multiple goroutines read/write it without synchronization.
The fix: pass by value (function argument), use a channel, or use sync.

**Compare with `[]byte` — why slices are dangerous to share:**

```go
  data := []byte("hello")
  go func() {
      data[0] = 'H'    // ← DATA RACE! Modifying shared backing array
  }()
  fmt.Println(data)     // ← DATA RACE! Reading while goroutine writes
```

Slices are mutable → shared backing array → data race.
Strings are immutable → shared backing array → perfectly safe.

This is Go's design: **strings trade mutation convenience for concurrency safety**.
When you need to mutate, convert to `[]byte`, do your work, convert back to `string`.
The conversion copies the bytes — creating a new, independent mutable buffer.

### The Cost of Immutability

Every "modification" creates a **new string** with a **new backing array**:

```go
s := "hello"
s = s + " world"   // allocates new 11-byte array, copies "hello" + " world"
                    // old "hello" backing bytes become GC-eligible
```

This is why string concatenation in loops is expensive — each `+=` allocates.

---

## 3. Bytes vs Runes vs Characters

This is the most important mental model for strings in Go. There are **three levels**:

```
"Hello, 世界"

BYTES (what the computer sees):
┌────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┐
│ 48 │ 65 │ 6C │ 6C │ 6F │ 2C │ 20 │ E4 │ B8 │ 96 │ E7 │ 95 │ 8C │
└────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┘
  H    e    l    l    o    ,  (sp)  ├── 世 ──┤  ├── 界 ──┤
                                    3 bytes     3 bytes
len("Hello, 世界") = 13 bytes

RUNES (Unicode code points — what Go uses):
┌─────┬─────┬─────┬─────┬─────┬─────┬─────┬───────┬───────┐
│ 'H' │ 'e' │ 'l' │ 'l' │ 'o' │ ',' │ ' ' │  '世'  │  '界'  │
└─────┴─────┴─────┴─────┴─────┴─────┴─────┴───────┴───────┘
len([]rune("Hello, 世界")) = 9 runes

CHARACTERS (what humans see):
H  e  l  l  o  ,     世  界     = 9 visible characters (in this case = runes)
```

### What Is a Rune?

```go
type rune = int32  // alias for int32 — holds any Unicode code point
```

A rune is a **Unicode code point** — a number that identifies a character in the Unicode
standard. `'世'` is rune `U+4E16` = `19990` in decimal.

### UTF-8 Encoding — Variable Width

Go strings are **UTF-8 encoded**. Each rune takes 1-4 bytes:

```
┌──────────────────────┬───────┬───────────────────────┐
│ Character Range      │ Bytes │ Examples              │
├──────────────────────┼───────┼───────────────────────┤
│ U+0000 - U+007F     │ 1     │ ASCII: A, z, 9, @     │
│ U+0080 - U+07FF     │ 2     │ é, ñ, Ω, ×            │
│ U+0800 - U+FFFF     │ 3     │ 世, 界, ☺, €           │
│ U+10000 - U+10FFFF  │ 4     │ 😀, 🎉, 𝕳              │
└──────────────────────┴───────┴───────────────────────┘
```

### The Implication

**Indexing by byte is NOT indexing by character:**

```go
s := "café"
s[3]           // 0xC3 — first byte of 'é', NOT the letter 'é'
len(s)         // 5 bytes (c=1, a=1, f=1, é=2)
len([]rune(s)) // 4 runes (c, a, f, é)
```

This is why you can't do `s[i]` and expect a character — you get a **byte**.

### When Characters ≠ Runes

Some visible characters are made of **multiple runes** (combining characters):

```go
s := "é"        // could be 1 rune (U+00E9) or 2 runes (e + U+0301 combining accent)
s = "👨‍👩‍👧"  // 5 runes: 👨 + ZWJ + 👩 + ZWJ + 👧 — but 1 visible "character"
```

For most server-side Go code, **rune = character** is close enough. If you need
true grapheme cluster counting (emoji, combining marks), use
`golang.org/x/text/unicode/norm` or third-party libraries.

---

## 4. String Iteration — Two Different Loops

### `for i` — Iterates Over Bytes

```go
s := "Go世界"
for i := 0; i < len(s); i++ {
    fmt.Printf("byte[%d] = %x\n", i, s[i])
}
// byte[0] = 47  (G)
// byte[1] = 6f  (o)
// byte[2] = e4  ┐
// byte[3] = b8  ├── 世 (3 bytes)
// byte[4] = 96  ┘
// byte[5] = e7  ┐
// byte[6] = 95  ├── 界 (3 bytes)
// byte[7] = 8c  ┘
```

### `for range` — Iterates Over Runes (UTF-8 decoded)

```go
s := "Go世界"
for i, r := range s {
    fmt.Printf("rune at byte %d = %c (U+%04X)\n", i, r, r)
}
// rune at byte 0 = G (U+0047)
// rune at byte 1 = o (U+006F)
// rune at byte 2 = 世 (U+4E16)   ← jumped from byte 2 to byte 5
// rune at byte 5 = 界 (U+754C)
```

**Key difference:** `for range` on a string automatically **decodes UTF-8**. The index
`i` is the **byte position**, and `r` is the **decoded rune**. Notice the index jumps
from 2 to 5 — it skipped the 3 bytes that encode `世`.

### What Happens with Invalid UTF-8?

```go
s := "ab\xff\xfec"  // invalid UTF-8 bytes
for _, r := range s {
    fmt.Printf("%c ", r)
}
// a b � � c
// Invalid bytes produce U+FFFD (replacement character)
```

`for range` is safe — it never panics on invalid UTF-8. It substitutes `U+FFFD`.

---

## 5. String ↔ Byte Slice Conversions — The Hidden Cost

### `[]byte(s)` — String to Byte Slice

```go
s := "hello"
b := []byte(s)   // allocates new 5-byte array, copies "hello" into it
```

**This allocates and copies.** The runtime must copy because:
- Strings are immutable — the backing bytes must not be modified
- Byte slices are mutable — the caller expects to modify `b`
- Sharing the same backing would violate string immutability

### `string(b)` — Byte Slice to String

```go
b := []byte{72, 101, 108, 108, 111}
s := string(b)   // allocates new backing array, copies bytes
```

**Also allocates and copies.** The runtime must copy because:
- If it shared the backing array, modifying `b[0]` later would change `s`
- That would violate string immutability

### The Cost

```go
// In a hot loop processing 50k requests/sec:
func handle(s string) {
    b := []byte(s)     // allocation + copy — 50k allocations/sec
    process(b)
    result := string(b) // another allocation + copy — 50k more
    return result
}
```

Each conversion is `O(n)` in the length of the string. In hot paths, this creates
significant GC pressure.

### Compiler Optimizations — When It Doesn't Copy

The compiler is smart enough to skip the copy in certain patterns:

```go
// 1. Map lookup with []byte key — no allocation
m := map[string]int{"hello": 1}
b := []byte("hello")
_ = m[string(b)]          // compiler optimizes: no copy, temporary string header

// 2. String comparison — no allocation
if string(b) == "hello" { ... }  // compiler: compare bytes directly, no alloc

// 3. Range over []byte as string — no allocation
for _, r := range string(b) { ... }  // compiler: decode UTF-8 in-place
```

### Zero-Copy Conversion (Go 1.20+ — `unsafe`)

Normal `[]byte → string` conversion **always copies** to preserve immutability:

```
[]byte b                            string s = string(b)

┌──────────────────┐               ┌──────────────────┐
│ ptr ──────────┐  │               │ ptr ──────────┐  │
│ len: 5        │  │               │ len: 5        │  │
└───────────────┼──┘               └───────────────┼──┘
                ↓                                  ↓
  Memory block A:                    Memory block B:     ← NEW allocation
  [h][e][l][l][o]                    [h][e][l][l][o]     ← bytes COPIED
```

Go must copy because `[]byte` is mutable — if they shared memory, modifying `b[0]`
would silently change the "immutable" string.

`unsafe.String` skips the copy entirely — the string header points directly at the
byte slice's backing array:

```go
import "unsafe"

b := []byte("hello")
s := unsafe.String(&b[0], len(b))  // NO copy, NO allocation
```

```
[]byte b                          string s

┌──────────────────┐             ┌──────────────────┐
│ ptr ──────────┐  │             │ ptr ──────────┐  │
│ len: 5        │  │             │ len: 5        │  │
└───────────────┼──┘             └───────────────┼──┘
                ↓                                ↓
                └───────── SAME memory ──────────┘
                  [h][e][l][l][o]
```

**The danger**: modifying the byte slice now corrupts the "immutable" string:

```go
b[0] = 'X'
fmt.Println(s)   // prints "Xello" — immutability VIOLATED
```

Use only in performance-critical code where you control the byte slice lifecycle
and can guarantee no further modifications (e.g., pprof shows `slicebytetostring`
allocations as a bottleneck at high RPS).

---

## 6. Substrings Share Backing Data — Memory Leak Trap

Substrings do **not** copy data. They create a new string header pointing into the
**same backing bytes**:

```go
huge := loadHugeString()   // 10 MB string
prefix := huge[:10]         // new header: {ptr: same as huge, len: 10}
```

```
huge header:    { ptr ──────────▶ [10 MB of bytes...............] }
                                   ▲
prefix header:  { ptr ─────────────┘, len: 10 }
```

The prefix is only 10 bytes, but the **entire 10 MB backing array** stays alive
because `prefix` holds a pointer into it. GC cannot collect the 10 MB until both
`huge` and `prefix` are dead.

### The Fix

```go
// Go 1.20+
prefix := strings.Clone(huge[:10])  // allocates new 10-byte backing array

// Before Go 1.20
prefix := string([]byte(huge[:10])) // string → []byte (copies 10) → string (copies 10)
```

`strings.Clone` explicitly allocates a new backing array, breaking the reference
to the original 10 MB.

### When This Matters

- Parsing large HTTP response bodies and keeping small substrings
- Extracting fields from large log lines
- Any pattern where a small substring outlives a large parent string

---

## 7. String Concatenation — Performance Analysis

### `+` Operator — Simple But Costly

```go
s := "hello"
s = s + " " + "world"  // allocates new 11-byte array, copies everything
```

Each `+` allocates a new backing array. In a loop, this is **O(n²)**:

```go
// ❌ O(n²) — each iteration allocates and copies everything so far
var s string
for i := 0; i < 1000; i++ {
    s += "x"  // iteration 1: copy 1 byte, iteration 2: copy 2, ... iteration 1000: copy 1000
}
// Total bytes copied: 1 + 2 + 3 + ... + 1000 = 500,500 copies + 1000 allocations
```

### `strings.Builder` — The Right Way

```go
// ✅ O(n) — amortized, like slice append
var b strings.Builder
b.Grow(1000)  // optional: pre-allocate if you know the size
for i := 0; i < 1000; i++ {
    b.WriteString("x")
}
s := b.String()
// Total: ~1-3 allocations (grows like a slice), 1000 bytes copied once each
```

`strings.Builder` internally uses a `[]byte` and calls `append`. It grows with the
same strategy as slices (2x, then ~1.25x). The final `String()` call uses
`unsafe.String` to convert the `[]byte` to a `string` **without copying** — this is
safe because Builder ensures the bytes won't be modified after `String()` is called.

#### How Builder Uses `unsafe.String` Safely

The full lifecycle, step by step:

```
Step 1: Create             var b strings.Builder        (buf is nil)
Step 2: Write              b.WriteString("hello")       (allocates buf, appends)
Step 3: Write more         b.WriteString(" world")      (appends to buf)

  Builder b
  ┌──────────────────────┐
  │ buf:                 │
  │   ptr ─────────┐     │
  │   len: 11      │     │
  │   cap: 16      │     │
  └────────────────┼─────┘
                   ↓
    [h][e][l][l][o][ ][w][o][r][l][d][_][_][_][_][_]
     ←──── len: 11 ────→                ←─ spare ──→

Step 4: Get string         s := b.String()

  Builder b.buf              returned string s
      │                           │
      │   ptr ────────┐           │  ptr ────────┐
      │   len: 11     │           │  len: 11     │
      └───────────────┼───        └──────────────┼──
                      ↓                          ↓
                      └──── SAME memory ─────────┘
                      [h][e][l][l][o][ ][w][o][r][l][d]

  Zero copy. Zero allocation. s shares buf's backing array.

Step 5: Done               Don't write to b anymore — s is safe
```

#### Builder's Copy Protection — The `addr` Trick

Builder stores a pointer to **itself** to detect if someone copies the struct:

```go
// Simplified from Go source
type Builder struct {
    addr *Builder   // points to itself after first use
    buf  []byte
}

func (b *Builder) copyCheck() {
    if b.addr == nil {
        b.addr = b               // first use: "I live at 0xABC"
    } else if b.addr != b {
        panic("strings: illegal use of non-zero Builder copied by value")
    }
}
```

Why? If you copy a Builder, both copies share the same `buf` backing array:

```go
b1 := strings.Builder{}
b1.WriteString("hello")
b2 := b1                    // struct copy — b2.buf points to same array!
b2.WriteString("DANGER")    // would corrupt b1's data → PANICS instead
```

```
b1                                b2 (COPY)
┌──────────────────┐             ┌──────────────────┐
│ addr: &b1        │             │ addr: &b1 ← stale│
│ buf → ───────┐   │             │ buf → ───────┐   │
└──────────────┼───┘             └──────────────┼───┘
               ↓                                ↓
               └────── SAME backing array ──────┘

copyCheck: &b2 != b2.addr (&b1) → PANIC
```

Every write method calls `copyCheck()`. Copied Builder detected → panic.

> **Design insight**: This is classic Go pragmatism. Rust would prevent this at
> compile time with ownership rules. Go says: "panic at runtime, the name is
> `unsafe`, you've been warned." Simple, clear, your responsibility.

### `fmt.Sprintf` — Convenient But Slow

```go
s := fmt.Sprintf("%s %s %d", first, last, age)
// Uses reflection internally to inspect types, allocates multiple intermediate values
```

`fmt.Sprintf` uses `reflect` under the hood to determine types. In hot paths, prefer
`strings.Builder` or `strconv` for simple conversions.

### `strings.Join` — Best for Known Slices

```go
parts := []string{"hello", "world", "go"}
s := strings.Join(parts, ", ")  // "hello, world, go"
// Pre-calculates total length → single allocation → copies each part once
```

### Concatenation Cost Summary

```
┌──────────────────────┬──────────────┬──────────────────────────────────┐
│ Method               │ Allocations  │ When to Use                      │
├──────────────────────┼──────────────┼──────────────────────────────────┤
│ s1 + s2              │ 1 per +      │ 2-3 strings, one-time            │
│ s += x (in loop)     │ N (O(n²))    │ Never in loops                   │
│ strings.Builder      │ ~1-3 (amort) │ Building strings incrementally   │
│ strings.Join         │ 1            │ Joining a known []string         │
│ fmt.Sprintf          │ Multiple     │ Formatted output, not hot paths  │
│ strconv.Itoa/AppendX │ 0-1          │ Number → string in hot paths     │
└──────────────────────┴──────────────┴──────────────────────────────────┘
```

---

## 8. String Comparison and Interning

### Comparison with `==`

String comparison in Go is **value-based** — it compares the actual bytes, not pointers:

```go
a := "hello"
b := string([]byte{'h', 'e', 'l', 'l', 'o'})
a == b   // true — same bytes, even though different backing arrays
```

**Under the hood:** The runtime first compares `len`. If lengths differ, returns `false`
immediately (O(1)). If lengths match, compares bytes via `memequal` (O(n)).

### String Interning

String literals in the same binary **may** share backing data. The compiler and linker
can deduplicate identical string constants:

```go
a := "hello"   // points to .rodata section
b := "hello"   // may point to SAME bytes in .rodata (compiler decides)
```

But this is a **compiler optimization**, not a language guarantee. Don't rely on
pointer equality for strings — always use `==`.

### Case-Insensitive Comparison

```go
// ✅ Correct for ASCII and Unicode
strings.EqualFold("Hello", "hello")   // true
strings.EqualFold("Ω", "ω")          // true (Greek omega)

// ❌ Wrong — only works for ASCII
strings.ToLower(a) == strings.ToLower(b)  // allocates two new strings just to compare
```

`EqualFold` compares without allocating — it decodes runes and folds case on the fly.

---

## 9. Strings and the Compiler — Escape Analysis

### String Literals Don't Allocate

```go
func greet() string {
    return "hello"   // no allocation — points directly to .rodata in the binary
}
```

String literals live in the binary's read-only data section. Returning a string literal
just copies the 16-byte header (pointer + length). The pointer points into the binary
itself — no heap allocation, no GC involvement.

### Conversions and Concatenation Force Allocation

```go
func process(b []byte) string {
    return string(b)        // allocates: must copy bytes to create immutable backing
}

func combine(a, b string) string {
    return a + b            // allocates: new backing array for concatenated result
}
```

### Passing Strings to `interface{}` — The Boxing Cost

```go
func log(msg string) {
    fmt.Println(msg)        // msg escapes to heap — Println takes any (interface)
}
```

Same principle as slices — boxing into an `eface` requires a heap pointer.
Every value passed to `fmt.Println(a ...any)` gets boxed:

```
Your int 42:                          Boxed into any (eface):
  just 8 bytes on stack               ┌──────────────────┐
  ┌─────┐                             │ type: *intType   │
  │  42 │               →             │ data: ptr ───────┼→ heap: [42]  ← allocation!
  └─────┘                             └──────────────────┘

Your string "hello":                  Boxed into any (eface):
  16 bytes on stack                    ┌──────────────────┐
  ┌──────────────┐                     │ type: *stringType│
  │ ptr → bytes  │       →            │ data: ptr ───────┼→ heap: {ptr, len}  ← alloc!
  │ len: 5       │                     └──────────────────┘
  └──────────────┘
```

At 10k RPS logging 5 fields each → 50k interface boxing allocations/sec → GC pressure.

This is not a theoretical concern. This is the exact problem that led Uber to build
`zap`, and eventually pushed the Go team to add `slog` to the standard library.

---

### 9a. From Strings to Production Logging — A Case Study

This section traces how understanding Go's string internals directly shaped two of
the most important logging libraries in the Go ecosystem. It's a real-world story
of how language internals affect production architecture.

#### The Problem: logrus at Uber (~2016)

Uber was running Go services at **millions of requests/sec**. Their logger was
`logrus` — the most popular Go logger at the time. Every log call looked like:

```go
logrus.WithFields(logrus.Fields{
    "user":    userID,        // string → boxed into interface{}
    "status":  200,           // int → boxed into interface{}
    "latency": 2.5,           // float64 → boxed into interface{}
}).Info("request completed")
```

`logrus.Fields` is `map[string]interface{}`. Every value gets boxed into an `eface`.
The map itself allocates. At Uber's scale:

```
Per log line:
  1 map allocation + N interface boxing allocations + map bucket allocations
  = ~5-10 heap allocations per log call

At 1M RPS × 5 log lines each = 5M log calls/sec
  = 25-50 MILLION allocations/sec just for logging
  = GC running constantly → p99 latency spikes every few seconds
```

Uber's engineers traced their GC-induced latency spikes directly to logging
allocations. They needed a logger with **zero heap allocations** in the hot path.

#### The Solution: zap — Zero Allocation Producer (2016)

Uber built `zap` from scratch with one goal: **eliminate every allocation**.

**Key Design Decision 1 — Typed Fields, No interface{}**

```go
// zap's Field struct — stores values directly, no boxing
type Field struct {
    Key       string       // 16 bytes
    Type      FieldType    // 8 bytes (enum tag)
    Integer   int64        // 8 bytes — numerics stored HERE
    String    string       // 16 bytes — strings stored HERE
    Interface interface{}  // 16 bytes — only for complex objects
}
// Total: ~64 bytes per field
```

The critical insight: instead of one `any` field that boxes everything, zap has
**dedicated typed fields**. Each value type goes to its specific slot:

```go
zap.String("user", "mert")     // → Field{Key:"user", Type:StringType, String:"mert"}
zap.Int("status", 200)         // → Field{Key:"status", Type:Int64Type, Integer:200}
zap.Float64("latency", 2.5)    // → Field{Key:"latency", Type:Float64Type, Integer:bits}
zap.Bool("cached", true)       // → Field{Key:"cached", Type:BoolType, Integer:1}
```

```
logrus (map[string]interface{}):        zap (typed Field struct):

  int 200 → box into eface:              int 200 → store directly:
  ┌──────────────┐                       ┌──────────────────────┐
  │ type: *int   │                       │ Type:    Int64Type   │
  │ data: ptr ───┼→ heap: [200]          │ Integer: 200         │  ← NO heap!
  └──────────────┘  ↑ allocation!        └──────────────────────┘

  string "mert" → box into eface:        string "mert" → direct field:
  ┌──────────────┐                       ┌──────────────────────┐
  │ type: *str   │                       │ Type:   StringType   │
  │ data: ptr ───┼→ heap: {ptr,len}      │ String: "mert"       │  ← NO heap!
  └──────────────┘  ↑ allocation!        └──────────────────────┘
```

Notice: zap's `Field` has a dedicated `String string` field (16 bytes) that stores
string headers **directly** — no eface boxing needed even for strings. This costs
an extra 16 bytes per Field struct, but eliminates the last allocation.

**Key Design Decision 2 — Encoder Writes Directly to Buffer**

Most loggers: build a `map` → serialize with `json.Marshal` → write string.
zap: write directly to a byte buffer, no intermediate representations:

```go
// Simplified zap JSON encoder — no reflection, no json.Marshal
func (enc *jsonEncoder) AddString(key, val string) {
    enc.buf.AppendByte('"')
    enc.buf.AppendString(key)       // direct byte copy
    enc.buf.AppendString(`":"`)
    enc.buf.AppendString(val)       // direct byte copy
    enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) AddInt64(key string, val int64) {
    enc.buf.AppendByte('"')
    enc.buf.AppendString(key)
    enc.buf.AppendString(`":`)
    enc.buf.AppendInt(val)          // strconv-style, no Sprintf
}
```

No `fmt.Sprintf`. No `json.Marshal`. No reflection. Hand-written byte operations.

**Key Design Decision 3 — Buffer Pooling with sync.Pool**

```go
var _pool = sync.Pool{
    New: func() interface{} {
        return &Buffer{bs: make([]byte, 0, 1024)}
    },
}

// Each log call:
// 1. Get buffer from pool     (no allocation if one available)
// 2. Write all fields to it   (direct byte manipulation)
// 3. Flush to output           (single write syscall)
// 4. Reset and return to pool  (ready for next call)
```

```
Request 1:  Pool → [buffer] → encode → flush → reset → [buffer] → Pool
Request 2:  Pool → [buffer] → encode → flush → reset → [buffer] → Pool
                    same buffer reused! zero allocations!
```

**Key Design Decision 4 — Level Check First**

```go
// Most loggers: build fields first, THEN check level → wasted work
logger.Debug("query", zap.String("sql", buildExpensiveSQL()))
// Even if Debug is disabled: buildExpensiveSQL() runs, Field struct built

// zap's Check pattern for extreme optimization:
if ce := logger.Check(zap.DebugLevel, "query"); ce != nil {
    ce.Write(zap.String("sql", buildExpensiveSQL()))
}
// Debug disabled: Check returns nil, nothing else executes
```

**The Full Zero-Alloc Pipeline:**

```
zap.Info("request",
    zap.String("method", "GET"),
    zap.Int("status", 200),
    zap.Duration("latency", d),
)

  1. Level check:  Info ≥ configured level? → yes, continue
  2. Fields:       [3]Field built on stack (no map, no interface boxing)
  3. Buffer:       get pre-allocated buffer from sync.Pool
  4. Encode:       write JSON directly to buffer bytes (no reflection)
  5. Output:       single write() syscall to stdout/file
  6. Cleanup:      reset buffer, return to pool

  Total heap allocations: 0
```

#### The Sequel: slog Enters the Standard Library (Go 1.21, 2023)

Seven years after zap, the Go team added `slog` to the stdlib. They studied zap's
approach and made a **different tradeoff**:

```go
// slog's Value struct — optimized for SIZE over zero-alloc
type Value struct {
    any  any       // strings and complex types go here (boxing needed)
    num  uint64    // numerics stored directly via bit-casting
    kind Kind      // type tag
}
// Total: ~32 bytes (half of zap's 64)
```

slog chose a **smaller struct** at the cost of **one allocation for strings**:

```go
slog.Int("status", 200)         // → Value{num: 200, kind: KindInt64}     → 0 allocs ✅
slog.Float64("lat", 3.14)       // → Value{num: Float64bits(3.14)}        → 0 allocs ✅
slog.Bool("cached", true)       // → Value{num: 1, kind: KindBool}        → 0 allocs ✅
slog.String("user", "mert")     // → Value{any: "mert", kind: KindString} → 1 alloc  ⚠️
```

Why does the string box? A string header is 16 bytes `{ptr, len}`. The `num` field
is `uint64` — only 8 bytes. 16 doesn't fit in 8. The Go team chose to keep `Value`
at 32 bytes rather than adding a second `uint64` field (which would make it 40 bytes
for every type, even ints and bools).

#### The Tradeoff Matrix

```
                          logrus         zap              slog
  ────────────────────    ──────────     ──────────       ──────────
  Year                    2014           2016             2023
  In stdlib?              ❌              ❌                ✅
  Field storage           map[string]any Field struct     Attr struct
  String allocs           1 per field    0                1 per string
  Numeric allocs          1 per field    0                0
  Field struct size       N/A (map)      ~64 bytes        ~32 bytes
  Serialization           json.Marshal   hand-written     Handler interface
  Buffer reuse            ❌              sync.Pool        Handler-dependent
  Reflection              ✅              ❌                ❌ in hot path
  
  Allocs per log line (5 fields):
    logrus:   ~5-10       (map + boxing + serialization)
    zap:      ~0          (everything on stack + pooled buffer)
    slog:     ~2          (string boxing only)
```

#### Why This Story Matters

This is a case study of how **one language design decision** — using `interface{}`
for generic parameters — ripples through the entire ecosystem:

```
  string is 16 bytes (header)
      ↓
  Can't fit in interface{} data pointer without heap allocation
      ↓
  Every log call that takes interface{} allocates
      ↓
  At Uber's scale: millions of allocs/sec from logging alone
      ↓
  GC pressure → p99 latency spikes
      ↓
  Uber builds zap: typed fields, zero allocs, buffer pooling
      ↓
  Go team studies zap → adds slog to stdlib with similar pattern
      ↓
  Today: typed methods (slog.String, zap.String) are the standard
         approach for allocation-sensitive APIs in Go
```

Understanding string internals isn't academic — it directly explains why production
Go code looks the way it does. The typed-method pattern you see in `slog` and `zap`
exists specifically because someone understood that a 16-byte string header can't fit
in an 8-byte eface data field without allocating.

#### When to Use Which

- **slog**: New projects, stdlib-only preference, most services under ~50k RPS.
  It's "good enough" for the vast majority of production workloads
- **zap**: Extreme throughput (>100k RPS), every allocation matters, already in your
  dependency tree. Battle-tested at Uber, Cloudflare, and similar scale
- **zerolog**: Another zero-alloc logger (different approach — method chaining with
  a builder pattern). Used by some high-performance projects
- **logrus**: Legacy. Still works, but not recommended for new high-throughput services

> **Go Wisdom**: *"A little copying is better than a little dependency."*
> But at Uber's scale, the allocation cost of `interface{}` was not "a little" — it
> was millions of heap allocations per second. Sometimes the dependency is worth it.

---

## 10. Production Patterns and Best Practices

### Pattern 1: Use `[]byte` Internally, `string` at API Boundaries

```go
// Internal processing: work with []byte to avoid conversion costs
func processRequest(body []byte) []byte {
    // parse, transform, build response — all as []byte
    return result
}

// API boundary: accept and return string
func HandleRequest(input string) string {
    b := []byte(input)          // one conversion in
    result := processRequest(b)
    return string(result)        // one conversion out
}
```

### Pattern 2: Pre-Size `strings.Builder`

```go
func buildResponse(items []Item) string {
    // Estimate: ~50 bytes per item
    var b strings.Builder
    b.Grow(len(items) * 50)    // single allocation
    for _, item := range items {
        b.WriteString(item.Name)
        b.WriteByte(',')
    }
    return b.String()
}
```

### Pattern 3: Avoid `fmt.Sprintf` in Hot Paths

```go
// ❌ Uses reflection, multiple allocations
key := fmt.Sprintf("user:%d:session:%s", userID, sessionID)

// ✅ Zero-reflection, minimal allocations
var b strings.Builder
b.WriteString("user:")
b.WriteString(strconv.Itoa(userID))
b.WriteString(":session:")
b.WriteString(sessionID)
key := b.String()
```

### Pattern 4: `strings.Clone` for Long-Lived Substrings

```go
func extractToken(header string) string {
    // header is "Bearer eyJhbGci..." (could be large)
    token := header[7:]
    return strings.Clone(token)  // detach from header's backing bytes
}
```

### Pattern 5: Use `strconv.AppendX` for Zero-Alloc Number Formatting

```go
// Append directly to an existing []byte — no intermediate string
buf := make([]byte, 0, 64)
buf = strconv.AppendInt(buf, 42, 10)        // "42"
buf = strconv.AppendFloat(buf, 3.14, 'f', 2, 64)  // "3.14"
```

---

## 11. Quick Reference Card

```
STRING HEADER (16 bytes on 64-bit)        Source: runtime/string.go
──────────────────────────────────
  str unsafe.Pointer  →  backing byte array (immutable)
  len int             →  byte count (NOT character count)

KEY RULES
─────────
  len(s)              →  byte count, O(1)
  len([]rune(s))      →  character count, O(n) — allocates rune slice
  utf8.RuneCountInString(s)  →  character count, O(n) — no allocation
  s[i]                →  byte at index i (NOT character)
  for i, r := range s →  iterate runes (UTF-8 decoded)
  for i := 0; i < len(s) →  iterate bytes (raw)

IMMUTABILITY
────────────
  s[0] = 'X'          →  compile error
  s + " world"        →  new string, new backing array
  s[:5]               →  new header, SHARED backing bytes

CONVERSIONS (both allocate + copy unless compiler optimizes)
────────────
  []byte(s)           →  mutable copy of string bytes
  string(b)           →  immutable string from byte slice
  []rune(s)           →  slice of Unicode code points
  string(r)           →  string from single rune

PERFORMANCE
───────────
  s += x (in loop)    →  O(n²) — never do this
  strings.Builder     →  O(n) amortized — use for building strings
  strings.Join        →  single allocation — use for joining slices
  strings.Clone       →  detach substring from parent (Go 1.20+)
  strconv.Itoa        →  int to string (no reflection)
  strconv.AppendInt   →  append int to []byte (zero alloc)

TOOLS
─────
  go build -gcflags='-m' ./...   →  see string escape analysis
  go test -bench=. -benchmem     →  measure allocations per operation
```

---

## One-Line Summary

> A string is a 16-byte immutable `{pointer, len}` header — essentially a read-only
> slice without capacity. Every "modification" creates a new backing array. `len()`
> counts bytes not characters. Use `strings.Builder` for concatenation, `[]rune` for
> character access, and `strings.Clone` to detach substrings from large parents.

---

## Further Reading

- [runtime/string.go](https://cs.opensource.google/go/go/+/master:src/runtime/string.go) — Runtime source for `stringStruct`, string-to-byte conversions, and concatenation internals
- [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings) — Official blog post on the relationship between bytes, runes, and UTF-8 encoding in Go
- [Go Spec — String types](https://go.dev/ref/spec#String_types) — Language specification for string semantics, indexing, and immutability guarantees
- [strings package source](https://cs.opensource.google/go/go/+/master:src/strings/builder.go) — `strings.Builder` implementation showing the `unsafe.String` zero-copy trick
- [Go Spec — Conversions to and from a string type](https://go.dev/ref/spec#Conversions_to_and_from_a_string_type) — Spec rules for `string([]byte)`, `[]byte(string)`, and `[]rune(string)` conversions
