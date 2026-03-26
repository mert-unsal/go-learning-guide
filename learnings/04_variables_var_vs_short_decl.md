# 14 — Variable Declaration: `var` vs `:=` vs Explicit Type

> When to use which form, what the compiler does, and the idiomatic conventions.

---

## The Three Forms

```go
var x int = 5     // Form 1: var + explicit type + value
var x = 5         // Form 2: var + inferred type
x := 5            // Form 3: short declaration + inferred type
```

**All three compile to identical machine code.** The difference is purely about
readability, context, and what the compiler infers.

---

## Form-by-Form Breakdown

### Form 1: `var x Type = value` — Explicit Type

```go
var count int = 0
var name string = "Go"
var r io.Reader = &MyStruct{}
```

**When to use:**
- When the **inferred type is not what you want**:

```go
var r io.Reader = &MyStruct{}   // you want io.Reader, not *MyStruct
var n int64 = 42                // you want int64, not int
var f float64 = 1               // you want float64, not int
```

- When making the type explicit **improves readability** in complex code

**Under the hood:** The compiler skips inference — it uses the type you declared
and checks that the value is assignable to it.

### Form 2: `var x = value` — Inferred Type

```go
var maxRetries = 3
var appName = "myservice"
var timeout = 30 * time.Second
```

**When to use:**
- **Package-level variables** (`:=` is not allowed at package level):

```go
package main

var Version = "1.0.0"   // ✅ package level
Version := "1.0.0"      // ❌ compile error: non-declaration statement outside function body
```

- When you want inference but `:=` isn't available

**Under the hood:** The compiler infers the type from the right-hand side, same as `:=`.

### Form 3: `x := value` — Short Declaration

```go
count := 0
name := "Go"
result := computeSomething()
```

**When to use:**
- **Inside functions — this is the default Go style**
- Most variable declarations in Go code use this form
- It's concise and the type is obvious from the value

**Under the hood:** Syntactic sugar for `var x = value`. Identical compiled output.

---

## The Zero Value Form: `var x Type` (No Assignment)

```go
var count int        // 0
var name string      // ""
var ok bool          // false
var s []int          // nil
var m map[string]int // nil
var err error        // nil
```

**This is the ONLY form where `var` has no `:=` equivalent.** You cannot write:

```go
count :=       // ❌ what value? Compile error.
```

**When to use:**
- When you want the **zero value** and assigning later
- Common pattern: declare before `if/switch`, assign inside, use after:

```go
var result int
if condition {
    result = computeA()
} else {
    result = computeB()
}
// use result here
```

---

## Where Each Form Is Allowed

```
┌─────────────────────┬─────────────────┬──────────────────┐
│ Form                │ Inside Function │ Package Level    │
├─────────────────────┼─────────────────┼──────────────────┤
│ var x Type = value  │ ✅              │ ✅               │
│ var x = value       │ ✅              │ ✅               │
│ x := value          │ ✅              │ ❌ not allowed   │
│ var x Type          │ ✅              │ ✅               │
└─────────────────────┴─────────────────┴──────────────────┘
```

**Key rule:** `:=` only works inside function bodies. At package level, you must use `var`.

---

## Type Inference: What the Compiler Decides

When you omit the type, the compiler infers from the value:

```go
x := 42              // int (not int32, not int64 — platform-dependent int)
x := 42.0            // float64 (not float32)
x := "hello"         // string
x := true            // bool
x := 'A'             // rune (which is int32)
x := 3 + 4i          // complex128

x := []int{1, 2, 3}  // []int
x := map[string]int{} // map[string]int
x := &User{}          // *User
```

**Gotcha — numeric defaults:**

```go
x := 42       // int — NOT int64
y := 42.0     // float64 — NOT float32

// If you need a specific numeric type, use explicit form:
var x int64 = 42
var y float32 = 42.0
```

---

## The `:=` Redeclaration Rule

`:=` can **redeclare** a variable if at least one variable on the left is new:

```go
x := 5
x := 10           // ❌ compile error: no new variables on left side

x, y := 5, "hi"
x, z := 10, true  // ✅ z is new, so x is reassigned (not redeclared)
```

This is commonly seen with `err`:

```go
result, err := doFirst()     // err declared
output, err := doSecond()    // err reassigned (output is new)
```

**Gotcha — shadowing in inner scopes:**

```go
x := 5
if true {
    x := 10    // ← this is a NEW x, shadows the outer x
    fmt.Println(x)  // 10
}
fmt.Println(x)      // 5 — outer x unchanged!

// Fix: use = (assignment) instead of := (declaration)
if true {
    x = 10     // ← assigns to the outer x
}
```

---

## Grouped `var` Declarations

At package level, group related variables:

```go
var (
    MaxRetries = 3
    Timeout    = 30 * time.Second
    AppName    = "myservice"
)
```

This is also used for compile-time interface checks:

```go
var _ io.Reader = (*MyStruct)(nil)  // compile error if MyStruct doesn't implement io.Reader
```

---

## The Idiomatic Rules

```
┌───────────────────────────────────────────────────────────────────┐
│  1. Inside functions → use :=  (default, most common)            │
│  2. Zero values → use var x Type  (no := equivalent)             │
│  3. Package level → use var  (:= not allowed)                    │
│  4. Explicit interface type → use var r InterfaceType = value     │
│  5. Explicit numeric type → use var n int64 = 42                 │
│  6. Group related package vars → use var ( ... ) block           │
└───────────────────────────────────────────────────────────────────┘
```

---

## Quick Reference

```
DECLARATION FORMS
─────────────────
  var x int = 5        explicit type, explicit value
  var x = 5            inferred type, explicit value
  x := 5               inferred type, explicit value (function only)
  var x int            explicit type, zero value

COMPILER BEHAVIOR
─────────────────
  All forms with same type+value → identical compiled output
  := is syntactic sugar for var x = value
  Type inference: int, float64, string, bool, rune for literals

SCOPE RULES
───────────
  := creates a NEW variable in current scope (can shadow outer)
  = assigns to EXISTING variable (no shadowing)
  := with multiple returns: redeclares if at least one var is new
```

---

## One-Line Summary

> `var` and `:=` produce identical compiled code — the choice is about **context**
> (package vs function level), **readability** (explicit vs inferred type), and
> **zero values** (only `var x Type` can declare without assigning).
