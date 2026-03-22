# Deep Dive: Go Interface Internals — iface, eface, itab & the Nil Trap

> Everything the runtime does when you assign a value to an interface,
> call a method through it, or compare it to nil.

---

## Table of Contents

1. [The Two Runtime Structs: `iface` and `eface`](#1-the-two-runtime-structs-iface-and-eface)
2. [The `itab` — The Method Dispatch Table](#2-the-itab--the-method-dispatch-table)
3. [Compile Time vs Runtime — Who Does What?](#3-compile-time-vs-runtime--who-does-what)
4. [Step-by-Step: Interface Assignment](#4-step-by-step-interface-assignment)
5. [Step-by-Step: Interface Method Call](#5-step-by-step-interface-method-call)
6. [The Three States of an Interface](#6-the-three-states-of-an-interface)
7. [The Nil Interface Trap](#7-the-nil-interface-trap)
8. [The Three Guards Against the Nil Trap](#8-the-three-guards-against-the-nil-trap)
9. [Type Assertions Under the Hood](#9-type-assertions-under-the-hood)
10. [Performance Implications](#10-performance-implications)
11. [Quick Reference Card](#11-quick-reference-card)

---

## 1. The Two Runtime Structs: `iface` and `eface`

Go uses **two different structs** for interface values at runtime, depending on whether
the interface has methods or not.

### `iface` — Non-Empty Interfaces (has methods)

Used for interfaces like `io.Reader`, `fmt.Stringer`, `error`, or any custom interface
with at least one method.

```
runtime.iface (16 bytes on 64-bit)
┌──────────────────────┬──────────────────────┐
│  tab  *itab          │  data unsafe.Pointer │
│  (method dispatch    │  (pointer to the     │
│   table + type info) │   concrete value)    │
└──────────────────────┴──────────────────────┘
```

- **`tab`** — pointer to an `itab` struct (see Section 2). Contains type metadata AND
  a cached array of method function pointers for this specific (interface, concrete type) pair.
- **`data`** — pointer to the actual value. If you assigned `&User{Name: "Bob"}`, this
  points to that `User` on the heap.

**Source:** `runtime/runtime2.go`
```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
```

### `eface` — Empty Interface (`interface{}` / `any`)

Used when the interface has **zero methods**. Since there's no method table needed,
Go uses a simpler, lighter struct.

```
runtime.eface (16 bytes on 64-bit)
┌──────────────────────┬──────────────────────┐
│  _type  *_type       │  data unsafe.Pointer │
│  (type descriptor    │  (pointer to the     │
│   only, no methods)  │   concrete value)    │
└──────────────────────┴──────────────────────┘
```

- **`_type`** — pointer to a `_type` struct (type descriptor: size, kind, hash, GC bitmap).
  No method table — just "what type is this?"
- **`data`** — same as `iface`, points to the actual value.

**Source:** `runtime/runtime2.go`
```go
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

### Visual Comparison

```
  Non-empty interface (e.g., io.Writer)       Empty interface (interface{} / any)
  ─────────────────────────────────────       ─────────────────────────────────
  ┌────────┬────────┐                         ┌────────┬────────┐
  │  tab ──┼──► itab│                         │ _type ─┼──► _type struct
  │        │  {     │                         │        │   {size, kind,
  │        │   inter│ ◄─ interface type info   │        │    hash, gcdata}
  │        │   _type│ ◄─ concrete type info    │        │
  │        │   fun[]│ ◄─ method pointers      │        │
  │        │  }     │                         │        │
  ├────────┼────────┤                         ├────────┼────────┤
  │ data ──┼──► concrete value on heap        │ data ──┼──► concrete value on heap
  └────────┴────────┘                         └────────┴────────┘
```

**Key insight:** Both are exactly 16 bytes (two pointers). The difference is what the
first pointer points TO — an `itab` (with methods) vs a bare `_type` (no methods).

---

## 2. The `itab` — The Method Dispatch Table

The `itab` is the heart of interface method dispatch. It answers: "For this specific
(interface type, concrete type) pair, where are the method implementations?"

```
runtime.itab
┌─────────────────────────────────────────────────────────────┐
│  inter  *interfacetype    // describes the interface        │
│         {                                                   │
│           typ    _type    // interface's own type info       │
│           pkgpath *string // package path                   │
│           mhdr   []imethod // list of methods the           │
│                            // interface requires            │
│         }                                                   │
├─────────────────────────────────────────────────────────────┤
│  _type  *_type            // describes the concrete type    │
│         {                                                   │
│           size    uintptr   // how many bytes                │
│           kind    uint8     // struct? ptr? int? etc.        │
│           hash    uint32    // for fast type comparison      │
│           str     nameOff   // type name                    │
│           ...               // GC bitmap, alignment, etc.   │
│         }                                                   │
├─────────────────────────────────────────────────────────────┤
│  hash   uint32            // copy of _type.hash for fast    │
│                           // type switch lookups            │
├─────────────────────────────────────────────────────────────┤
│  fun    [1]uintptr        // VARIABLE SIZE array of method  │
│         fun[0] = &(*EmailNotifier).Notify                   │
│         fun[1] = &(*EmailNotifier).OtherMethod              │
│         ...one entry per interface method, in sorted order  │
└─────────────────────────────────────────────────────────────┘
```

### How `itab` Gets Built and Cached

The runtime maintains a **global hash table** called `runtime.itabTable`.

```
                    itabTable (global, process-wide)
                    ────────────────────────────────
                    hash key = hash(inter, _type)

 ┌───────────────────────────────────────────────────────────────┐
 │ (io.Writer, *os.File)       → itab{ fun: [Write] }           │
 │ (io.Reader, *bytes.Buffer)  → itab{ fun: [Read] }            │
 │ (Notifier, *EmailNotifier)  → itab{ fun: [Notify] }          │
 │ (fmt.Stringer, *User)       → itab{ fun: [String] }          │
 │ ...                                                           │
 └───────────────────────────────────────────────────────────────┘
```

**Lifecycle:**

1. **First assignment** of `*EmailNotifier` to `Notifier` → runtime calls `runtime.getitab()`
2. `getitab()` checks `itabTable` for key `(Notifier, *EmailNotifier)`
3. **Cache miss** → runtime builds the `itab`:
   - Walks the interface's method list (`inter.mhdr`)
   - Walks the concrete type's method list (`_type.methods`)
   - Both are sorted alphabetically — matched with a **single O(n+m) merge walk**
   - Stores each matched method's function pointer in `fun[]`
4. **Stores** the new `itab` in `itabTable`
5. **All future assignments** of `*EmailNotifier` to `Notifier` → cache hit → O(1) lookup

**Important:** The `itab` is built **once per (interface, concrete type) pair** for the
entire lifetime of the process. It's never rebuilt or garbage collected.

---

## 3. Compile Time vs Runtime — Who Does What?

There is no `compileTime.iface`. The compiler and runtime have distinct, complementary roles:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          COMPILE TIME                                   │
│  (go build — the compiler, cmd/compile)                                │
│                                                                         │
│  1. TYPE CHECKING                                                       │
│     Verifies your concrete type has ALL methods the interface requires. │
│     If *EmailNotifier is missing Notify(), you get a compile error.     │
│     This happens at the assignment site:                                │
│       var n Notifier = &EmailNotifier{}  // compiler checks here       │
│                                                                         │
│  2. CODE GENERATION                                                     │
│     Emits assembly instructions to:                                     │
│     - Call runtime.getitab() to get/build the itab at runtime           │
│     - Construct the iface struct {tab, data}                            │
│     - For method calls: load tab → load fun[i] → indirect CALL         │
│                                                                         │
│  3. ESCAPE ANALYSIS                                                     │
│     Determines if the concrete value must escape to heap.               │
│     Assigning to an interface almost always forces a heap allocation    │
│     because iface.data is an unsafe.Pointer (must outlive the stack).   │
│                                                                         │
│  4. STATIC ITAB OPTIMIZATION (link time)                                │
│     For statically known assignments, the linker can pre-build some     │
│     itab entries — avoiding the runtime.getitab() call entirely.        │
├─────────────────────────────────────────────────────────────────────────┤
│                           RUNTIME                                       │
│  (the running binary — runtime package)                                │
│                                                                         │
│  1. ITAB TABLE                                                          │
│     Global hash table of cached itab entries.                           │
│     Built lazily on first use of each (interface, type) pair.           │
│                                                                         │
│  2. METHOD DISPATCH                                                     │
│     Indirect function call through itab.fun[i].                         │
│     CPU loads the function pointer and CALLs it.                        │
│                                                                         │
│  3. TYPE ASSERTIONS / TYPE SWITCHES                                     │
│     Compares itab._type.hash against target type hash.                  │
│     Fast path: hash match → verify → extract data.                     │
│                                                                         │
│  4. INTERFACE EQUALITY                                                  │
│     i == nil checks: tab == nil AND data == nil.                        │
│     i == j checks: same itab + deep-equal data.                        │
└─────────────────────────────────────────────────────────────────────────┘
```

**Verification tools:**
```bash
go build -gcflags='-m'     # escape analysis — what escapes to heap
go build -gcflags='-m -m'  # verbose escape analysis with reasons
go build -gcflags='-S'     # assembly output — see iface construction & dispatch
```

---

## 4. Step-by-Step: Interface Assignment

Let's trace exactly what happens when you write:

```go
type Notifier interface {
    Notify(msg string) error
}

type EmailNotifier struct{ smtp string }
func (e *EmailNotifier) Notify(msg string) error { /* ... */ }

var n Notifier = &EmailNotifier{smtp: "mx.google.com"}
```

### Step 1 — Compile Time: Type Check

```
Compiler sees: var n Notifier = &EmailNotifier{...}

  Does *EmailNotifier have Notify(string) error?
  ├─ Scan *EmailNotifier's method set
  ├─ Found: func (e *EmailNotifier) Notify(msg string) error  ✅
  └─ Type check passes — code generation proceeds
```

### Step 2 — Compile Time: Escape Analysis

```
Compiler decides: &EmailNotifier{smtp: "mx.google.com"}
  ├─ This value is assigned to an interface (iface.data = unsafe.Pointer)
  ├─ unsafe.Pointer could escape the current stack frame
  └─ Decision: ALLOCATE ON HEAP ← this is the cost of interfaces
```

### Step 3 — Runtime: Allocate the Value

```
  Heap:
  ┌──────────────────────────────┐
  │ EmailNotifier{               │
  │   smtp: "mx.google.com"     │  ◄── 0xc0000a4000 (heap address)
  │ }                            │
  └──────────────────────────────┘
```

### Step 4 — Runtime: Get or Build the `itab`

```
  runtime.getitab(Notifier, *EmailNotifier)
  ├─ Hash lookup in itabTable
  ├─ First time? Cache miss → BUILD:
  │   ├─ Walk Notifier.methods: [Notify(string) error]
  │   ├─ Walk (*EmailNotifier).methods: [Notify(string) error]
  │   ├─ Match: fun[0] = address of (*EmailNotifier).Notify
  │   └─ Store in itabTable
  └─ Return: *itab at 0x4a7b20
```

### Step 5 — Runtime: Assemble the `iface`

```
  n (on the stack):
  ┌──────────────────────┬──────────────────────┐
  │ tab = 0x4a7b20       │ data = 0xc0000a4000  │
  │ (itab for            │ (heap-allocated      │
  │  Notifier,           │  EmailNotifier)      │
  │  *EmailNotifier)     │                      │
  └──────────────────────┴──────────────────────┘
            │                       │
            ▼                       ▼
      itab{                   EmailNotifier{
        inter: Notifier         smtp: "mx.google.com"
        _type: *EmailNotifier }
        fun[0]: 0x4823a0  ──► (*EmailNotifier).Notify
      }
```

---

## 5. Step-by-Step: Interface Method Call

Now let's trace `n.Notify("server down")`:

```go
n.Notify("server down")
```

### Step 1 — Load the `itab`

```
  CPU reads n.tab → 0x4a7b20 (pointer to itab)
```

### Step 2 — Load the Method Pointer

```
  CPU reads itab.fun[0] → 0x4823a0 (address of (*EmailNotifier).Notify)

  fun[] index is determined at COMPILE TIME — the compiler knows Notify is
  method #0 in the interface's sorted method list. This is a constant offset.
```

### Step 3 — Load the Data Pointer

```
  CPU reads n.data → 0xc0000a4000 (pointer to EmailNotifier value)
```

### Step 4 — Indirect Call

```
  CALL 0x4823a0(0xc0000a4000, "server down")
       ▲              ▲              ▲
       │              │              └─ argument
       │              └─ receiver (e *EmailNotifier)
       └─ (*EmailNotifier).Notify function address

  This is an INDIRECT call — the CPU doesn't know the target until it
  reads the function pointer from the itab. This means:
  ├─ The compiler CANNOT inline this call
  ├─ The CPU branch predictor may mispredict on first call
  └─ Subsequent calls to the same (interface, type) pair are predicted correctly
```

### Visual Summary: The Full Call Chain

```
  n.Notify("server down")

  ┌─── n ────────────┐
  │ tab ──► itab     │     ┌─ EmailNotifier ─┐
  │ data ──┼─────────┼────►│ smtp: "mx..."   │
  └────────┼─────────┘     └─────────────────┘
           │                       ▲
           ▼                       │
     itab.fun[0] ──► (*EmailNotifier).Notify(receiver, "server down")
                                   │              ▲
                                   └──────────────┘
                                   receiver = n.data
```

---

## 6. The Three States of an Interface

Every interface value is in one of exactly three states:

```go
var s Stringer         // State 1: (nil, nil)
s = User{Name: "Bob"}  // State 2: (*itab, 0xc000...)
var u *User
s = u                  // State 3: (*itab, nil)   ← THE TRAP
```

### State 1 — True Nil Interface

```
var s Stringer

s:
┌──────────────┬──────────────┐
│     nil      │     nil      │
└──────────────┴──────────────┘

s == nil → TRUE ✅   (both words are zero)
```

Both `tab` and `data` are zero. This is what `nil` means for an interface.

### State 2 — Normal, Non-Nil Interface

```
s = User{Name: "Bob", Age: 25}

s:
┌──────────────┬──────────────┐
│  *itab(User) │  0xc000a080  │──► User{Name:"Bob", Age:25}
└──────────────┴──────────────┘

s == nil  → FALSE ✅
s.String() → works fine ✅
```

### State 3 — The Nil Pointer Trap ⚠️

```go
var u *User   // u is nil
s = u         // assign nil pointer to interface
```

```
s:
┌──────────────┬──────────────┐
│  *itab(User) │     nil      │
└──────────────┴──────────────┘

s == nil    → FALSE ❌  (tab is NOT nil — it has type info!)
s.String()  → PANIC ❌  (data IS nil — method dereferences nil receiver)
```

---

## 7. The Nil Interface Trap

### Why Does This Happen?

When Go assigns a value to an interface, it **always** writes the type into the `tab` field —
**regardless of whether the value itself is nil**.

```go
var u *User   // concrete nil pointer — the type (*User) IS known
s = u         // Go writes: tab = itab(*User), data = nil
              // it has no way to "forget" the type just because the value is nil
```

Go cannot collapse this back to `(nil, nil)` because the type information was known
statically at the assignment site. The compiler stamped `*User` into the `itab` and moved on.

### How `== nil` Works on Interfaces

The `==` operator on an interface checks:

```
(tab == nil) AND (data == nil)
```

Only when **both are zero** is the interface considered nil. State 3 has a non-nil `tab`,
so `== nil` is `false` — even though the data pointer inside is nil.

### The Classic Bug

```go
// BAD
func findUser(found bool) Stringer {
    var u *User              // type: *User, value: nil
    if !found {
        return u             // returns iface{tab: *itab, data: nil} — NOT nil!
    }
    return &User{Name: "Alice"}
}

result := findUser(false)
if result == nil {           // FALSE — the nil check is BYPASSED
    fmt.Println("not found")
}
result.String()              // PANIC — data is nil, method dereferences it
```

### What the Panic Trace Looks Like

```
n.Notify("hello")                    // interface method call
  ↓
runtime loads n.tab.fun[0]           → address of (*EmailNotifier).Notify ← SUCCEEDS
  ↓
calls (*EmailNotifier).Notify(nil, "hello")   ← receiver is nil
  ↓
inside Notify: e.smtpClient.Send()   → dereferences nil receiver 'e'
  ↓
PANIC: nil pointer dereference        ← error surfaces DEEP INSIDE the method

Stack trace points at EmailNotifier.Notify, NOT at where you assigned the typed nil.
This makes production debugging harder — the root cause is far from the crash site.
```

---

## 8. The Three Guards Against the Nil Trap

### Guard 1 — Fix at the Source (always prefer this)

**Never return a typed nil variable through an interface — return bare `nil` directly.**

```go
// BAD — typed nil leaks in
func findUser(found bool) Stringer {
    var u *User
    if !found {
        return u   // iface{*itab, nil} — bug!
    }
    return &User{Name: "Alice"}
}

// GOOD — return untyped nil
func findUser(found bool) Stringer {
    if !found {
        return nil  // iface{nil, nil} — true nil, == nil check works
    }
    return &User{Name: "Alice"}
}
```

Bare `return nil` in a function returning an interface produces `iface{nil, nil}`.

**When to use:** Every time you write a function returning an interface type. Make this a habit.

### Guard 2 — Type Assertion (when you know the concrete type)

```go
func safeCall(s Stringer) {
    u, ok := s.(*User)   // step 1: extract the concrete *User
    if !ok || u == nil {  // step 2: check type match AND nil pointer
        fmt.Println("skipping — nil or wrong type")
        return
    }
    fmt.Println(u.String()) // safe — u is a non-nil *User
}
```

Under the hood for `s.(*User)` when `s = iface{tab: *itab(*User), data: nil}`:

```
1. Compare s.tab._type.hash == hash(*User) → YES → ok = true
2. Extract s.data → nil → u = nil
∴ ok is true (type matched) but u is nil (pointer is nil)
  You must check BOTH.
```

**When to use:** When you know the expected concrete type (e.g., inside a service you control).

### Guard 3 — Reflect (last resort, for truly dynamic types)

```go
func isTrulyNil(i any) bool {
    if i == nil {
        return true  // fast path: iface{nil, nil}
    }
    v := reflect.ValueOf(i)
    switch v.Kind() {
    case reflect.Ptr, reflect.Interface,
         reflect.Slice, reflect.Map,
         reflect.Chan, reflect.Func:
        return v.IsNil()  // inspects the data word directly
    }
    return false  // int, string, struct etc. can never be nil
}
```

**When to use:** Generic libraries receiving `any` where you can't know the type. Avoid in hot paths
(reflection allocates and is ~10x slower than a type assertion).

### Decision Tree

```
You have an interface value and need to know if it's safe to call methods
│
├─ Can you change the code that PRODUCES the interface value?
│   └─ YES → Guard 1: return nil, not a typed nil variable
│
└─ NO, you receive it from outside
    │
    ├─ Do you know the concrete type?
    │   └─ YES → Guard 2: v, ok := i.(*Type); ok && v != nil
    │
    └─ NO, it's truly dynamic / any
        └─ Guard 3: reflect.ValueOf(i).IsNil()
```

---

## 9. Type Assertions Under the Hood

Given any interface variable `i` (e.g., `var i Notifier = &EmailNotifier{}`):

### Single Type Assertion: `v, ok := i.(ConcreteType)`

```
// i is any interface value — could be Notifier, io.Reader, error, any, etc.
// The assertion extracts the concrete type stored inside.

var i Notifier = &EmailNotifier{smtp: "mx.google.com"}
v, ok := i.(*EmailNotifier)   // is the concrete type inside i a *EmailNotifier?

Under the hood:
1. Read i.tab._type.hash                    // O(1) — hash stored in itab
2. Compare with hash(*EmailNotifier)        // O(1) — compiler knows this at build time
3. Match? → ok = true, v = *(*EmailNotifier)(i.data)   // cast the data pointer
   No match? → ok = false, v = nil (*EmailNotifier zero value)
```

### Type Switch: `switch v := i.(type) { ... }`

```go
// i is the interface variable you're inspecting.
// The switch tests which concrete type is stored inside.
var i any = 42

switch v := i.(type) {
case int:       // compare i._type.hash with hash(int)     ← eface, not itab (empty interface)
case string:    // compare i._type.hash with hash(string)
case *User:     // compare i._type.hash with hash(*User)
default:        // no match
}
```

The compiler generates a sequence of hash comparisons (or a jump table for many cases).
The hash is copied into `itab.hash` specifically to make this fast — no pointer chase needed.

---

## 10. Performance Implications

### Interface Method Calls Cannot Be Inlined

```go
func directCall(e *EmailNotifier) { e.Notify("hi") }  // ← CAN be inlined
func ifaceCall(n Notifier)        { n.Notify("hi") }  // ← CANNOT be inlined
```

Why? Inlining requires knowing the target function at compile time. With interfaces,
the target lives in `itab.fun[0]` which is only known at runtime.

**Production impact:** In hot paths processing 100k+ operations/sec, interface dispatch
overhead adds up. Consider using concrete types in inner loops and interfaces at boundaries.

### Interface Assignment Forces Heap Allocation

```go
var n Notifier = &EmailNotifier{smtp: "mx.google.com"}
//                ↑ this value ESCAPES to heap
//                  because iface.data is unsafe.Pointer
```

Every value assigned to an interface must live on the heap (the `data` pointer must remain
valid). This creates GC pressure. Verify with:

```bash
go build -gcflags='-m' ./...
# output: "EmailNotifier{} escapes to heap"
```

### Cost Summary

```
┌─────────────────────┬────────────────────┬─────────────────────────────────┐
│ Operation           │ Cost               │ Why                             │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Interface assign    │ ~30-50ns           │ itab lookup + possible heap     │
│ (first time)        │                    │ alloc + GC bookkeeping          │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Interface assign    │ ~5-15ns            │ itab cached, just heap alloc    │
│ (cached itab)       │                    │                                 │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Interface method    │ ~2-5ns overhead    │ Indirect call (load fun ptr +   │
│ call                │ vs direct call     │ CALL) — prevents inlining       │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ Type assertion      │ ~1-2ns             │ Hash comparison + pointer cast  │
├─────────────────────┼────────────────────┼─────────────────────────────────┤
│ == nil check        │ ~0.5ns             │ Two pointer comparisons         │
└─────────────────────┴────────────────────┴─────────────────────────────────┘
```

---

## 11. Quick Reference Card

```
NON-EMPTY INTERFACE (has methods)
─────────────────────────────────
runtime.iface { tab *itab, data unsafe.Pointer }
  └─ itab { inter *interfacetype, _type *_type, hash uint32, fun [N]uintptr }
     └─ fun[] = cached method pointers for (interface, concrete type) pair
     └─ cached globally in runtime.itabTable — built once, reused forever

EMPTY INTERFACE (interface{} / any)
───────────────────────────────────
runtime.eface { _type *_type, data unsafe.Pointer }
  └─ no method table, just type descriptor + data pointer
  └─ simpler because no methods to dispatch

NIL SEMANTICS
─────────────
  iface{nil, nil}    == nil  → TRUE   ← "true nil" interface
  iface{*itab, nil}  == nil  → FALSE  ← "typed nil" TRAP
  iface{*itab, ptr}  == nil  → FALSE  ← normal non-nil

COMPILE TIME vs RUNTIME
───────────────────────
  Compile: type checking, code generation, escape analysis, static itab optimization
  Runtime: itab caching, method dispatch (indirect call), type assertions, nil checks

TOOLS
─────
  go build -gcflags='-m'      # escape analysis: what goes to heap
  go build -gcflags='-m -m'   # verbose escape analysis with reasons
  go build -gcflags='-S'      # assembly output: see iface construction
  go test -race ./...         # race detector: always run
```

---

## One-Line Summary

> An interface is a two-word struct `{tab, data}` — `tab` points to the `itab`
> (method dispatch table cached per type pair), `data` points to the value.
> It's `nil` only when **both** words are zero. The compiler verifies types;
> the runtime dispatches methods. This is why interface calls can't be inlined.

