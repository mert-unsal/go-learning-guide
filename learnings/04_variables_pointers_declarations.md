# Deep Dive: Variables, Pointers & Declarations

> Declaration forms, type inference, auto-dereference, and addressability rules.

---

## Table of Contents

- [Part 1 — Variable Declarations: `var` vs `:=`](#part-1--variable-declarations-var-vs-)
- [Part 2 — Pointers: Auto-Dereference & Auto-Address](#part-2--pointers-auto-dereference--auto-address)
- [Quick Reference Card](#quick-reference-card)

---

## Part 1 — Variable Declarations: `var` vs `:=`

### The Three Forms

```go
var x int = 5     // Form 1: var + explicit type + value
var x = 5         // Form 2: var + inferred type
x := 5            // Form 3: short declaration + inferred type
```

**All three compile to identical machine code.** The difference is purely about
readability, context, and what the compiler infers.

---

### Form-by-Form Breakdown

#### Form 1: `var x Type = value` — Explicit Type

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

#### Form 2: `var x = value` — Inferred Type

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

#### Form 3: `x := value` — Short Declaration

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

### The Zero Value Form: `var x Type` (No Assignment)

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

### Where Each Form Is Allowed

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

### Type Inference: What the Compiler Decides

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

### The `:=` Redeclaration Rule

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

### Grouped `var` Declarations

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

### The Idiomatic Rules

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

## Part 2 — Pointers: Auto-Dereference & Auto-Address

Go's compiler does **two silent tricks** to make your code less verbose. Understanding exactly when and why they happen is the key.

### The Two Tricks

| Name | What Go does silently | Example |
|---|---|---|
| **Auto-address** | `c.Method()` → `(&c).Method()` | When method needs `*T` but you have `T` |
| **Auto-deref** | `p.Field` → `(*p).Field` | When accessing fields on a `*T` |

---

### Trick 1: Auto-Dereference (field & method access on pointers)

#### Rule: When you access a field or method through a pointer, Go silently inserts `*`

```go
type Person struct {
    Name string
    Age  int
}

p := &Person{Name: "Alice", Age: 30}

// You write:       p.Age++
// Go compiles it:  (*p).Age++
// Both are identical — Go does it for you
```

This works for **any depth of pointer**:

```go
p := &Person{Age: 30}
pp := &p   // pointer to a pointer

pp.Age++        // Go does: (*(*pp)).Age++  — all the way down
fmt.Println(pp) // still works
```

#### When does auto-deref happen?

✅ Always — whenever you use `.` (dot notation) on a pointer type.

```go
p.Name    // ✅ auto-deref → (*p).Name
p.Age++   // ✅ auto-deref → (*p).Age++
p.Method() // ✅ auto-deref if Method has value receiver
```

#### When must you write `*` explicitly?

❌ When you want the **whole value**, not a field:

```go
n := new(int)
*n = 10

fmt.Println(*n)   // ✅ must be explicit — you want the int value
fmt.Println(n)    // prints address: 0xc000018080

// Assigning the whole value:
*n = 99           // ✅ must be explicit — replacing the entire value
n = 99            // ❌ compile error — n is *int, not int
```

**Rule of thumb:** Use `*` explicitly when you're working with the **whole pointed-to value**, not a field/method inside it.

---

### Trick 2: Auto-Address (calling pointer-receiver methods on values)

#### Rule: When a method requires `*T` but you have `T`, Go silently inserts `&` — BUT only if the variable is addressable

```go
type Counter struct{ count int }

func (c *Counter) Increment() { c.count++ }  // pointer receiver

c := Counter{count: 0}  // c is a plain Counter (value, not pointer)

c.Increment()     // You write this
(&c).Increment()  // Go compiles it as this — c is addressable ✅
```

#### What is "addressable"?

A variable is addressable if Go can take its memory address with `&`.

```go
// ✅ ADDRESSABLE — auto-address works
c := Counter{}        // local variable
c.Increment()         // works — Go does (&c).Increment()

arr := [3]Counter{}
arr[0].Increment()    // works — array element is addressable

s := []Counter{{}}
s[0].Increment()      // works — slice element is addressable

// ❌ NOT ADDRESSABLE — auto-address FAILS, must be explicit
Counter{}.Increment()                  // ❌ compile error — temporary value, no address
getCounter().Increment()               // ❌ compile error — function return value, no address
map[string]Counter{}["x"].Increment()  // ❌ compile error — map values are not addressable
```

The temporary `Counter{}` has no memory address because Go hasn't assigned it a location yet. There's nothing to point to.

#### The fix for non-addressable cases:

```go
// Instead of:
Counter{}.Increment()          // ❌

// Do:
c := Counter{}
c.Increment()                   // ✅ now it has an address

// Or use a pointer from the start:
c := &Counter{}
c.Increment()                   // ✅ already a pointer
```

---

### The Full Decision Table

```
You have T,  method wants T  → direct call, nothing special
You have T,  method wants *T → Go auto-addresses: (&T).Method()  [if addressable]
You have *T, method wants *T → direct call, nothing special
You have *T, method wants T  → Go auto-derefs: (*T).Method()
```

---

### When to Write It Explicitly vs Leave It Implicit

#### Leave it implicit (let Go handle it):

```go
// Field access — always let Go handle it
p.Name = "Alice"       // NOT (*p).Name = "Alice"
p.Age++                // NOT (*p).Age++

// Method calls on variables — always let Go handle it
c.Increment()          // NOT (&c).Increment()
p.SomeMethod()         // NOT (*p).SomeMethod()
```

Writing `(*p).Age` or `(&c).Increment()` is valid Go, but it's **noisy and unusual** — no Go developer writes it that way in practice.

#### Write it explicitly when:

```go
// 1. You want the whole value behind a pointer
val := *p              // ✅ copying the entire struct
fmt.Println(*n)        // ✅ printing an int value, not its address

// 2. You're assigning to the whole pointed-to value
*n = 42                // ✅ replacing the entire value
*p = Person{"Bob", 25} // ✅ replacing the entire struct

// 3. You're taking an address to pass somewhere
someFunc(&x)           // ✅ passing address of x to a function
p := &Person{}         // ✅ creating a pointer to a new struct

// 4. Nil checks — always explicit
if p != nil {          // ✅ checking if pointer is nil
    fmt.Println(p.Name)
}

// 5. Double pointers (rare)
pp := &p               // ✅ pointer to a pointer
```

---

## Quick Reference Card

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

AUTO-DEREF  (p.Field, p.Method)
───────────
  → Go inserts * automatically
  → ALWAYS works on pointer types with dot notation
  → Write * explicitly only when grabbing/replacing the WHOLE value

AUTO-ADDRESS  (c.PointerMethod())
────────────
  → Go inserts & automatically
  → ONLY works when variable is ADDRESSABLE (has a memory slot)
  → Fails on: temporaries, function returns, map values
  → Write & explicitly when passing to functions or creating pointers

GOLDEN RULE
───────────
  Use dot notation normally (p.Name, c.Increment())
  Use * only when you mean "give me everything at this address"
  Use & only when you mean "give me the address of this thing"
```
