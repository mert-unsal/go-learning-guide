# Deep Dive: Interface Values, the Nil Trap, and How to Guard Against It

---

## What is an interface value, really?

When you think of an interface variable in Go, forget the idea of it being a simple
pointer or reference. At runtime, **every interface value is two words side by side in memory**:

```
┌──────────────┬──────────────┐
│  type word   │  value word  │
│  (itab ptr)  │  (data ptr)  │
└──────────────┴──────────────┘
```

- The **type word** points to metadata about the concrete type stored (methods, type name, size…).
- The **value word** points to the actual data (or stores it inline if it's small enough).

The Go runtime uses this pair for every single interface operation — method dispatch, equality checks, type assertions. Everything.

---

## The three states an interface can be in

```go
var s Stringer         // State 1: (nil, nil)
s = User{Name: "Bob"}  // State 2: (*User, 0xc000...)
var u *User
s = u                  // State 3: (*User, nil)   ← THE TRAP
```

Let's visualise each one:

### State 1 — True nil interface

```
var s Stringer

s:
┌──────────────┬──────────────┐
│     nil      │     nil      │
└──────────────┴──────────────┘

s == nil → TRUE ✅
```

Both words are zero. This is what `nil` means for an interface.

---

### State 2 — Normal, non-nil interface

```
s = User{Name: "Bob", Age: 25}

s:
┌──────────────┬──────────────┐
│  *itab(User) │  0xc000a080  │──► User{Name:"Bob", Age:25}
└──────────────┴──────────────┘

s == nil → FALSE ✅
s.String() → works fine ✅
```

The type word points to `User`'s type info. The value word points to the actual `User` data.

---

### State 3 — The nil pointer trap

```go
var u *User   // u is nil
s = u         // assign nil pointer to interface
```

```
s:
┌──────────────┬──────────────┐
│  *itab(User) │     nil      │
└──────────────┴──────────────┘

s == nil → FALSE ❌  (type word is NOT nil!)
s.String() → PANIC ❌ (value word IS nil, method tries to dereference it)
```

This is the trap. The moment you assign `u` to `s`, the **type word gets filled** with `*User`.
The interface is no longer `(nil, nil)` — it is `(*User, nil)`.
The `== nil` check only returns `true` when **both words are nil**, so it says `false`.
But calling any method will try to dereference that nil data pointer → **panic**.

---

## Why does this happen? The assignment rule.

When Go assigns a value to an interface, it always writes the type into the type word —
**regardless of whether the value itself is nil**.

```go
var u *User   // concrete nil pointer — type is known: *User
s = u         // Go writes: type=*User, value=nil
              // it has no way to "forget" the type just because the value is nil
```

Go cannot collapse this back to `(nil, nil)` automatically because the type information
was already known statically at the assignment site. The compiler stamped `*User` into
that type word and moved on.

---

## How `== nil` actually works on interfaces

The `==` operator on an interface checks:

```
(type word == nil) AND (value word == nil)
```

Only when **both are zero** is the interface considered nil.
State 3 has a non-zero type word, so `== nil` is `false` — even though the pointer inside is nil.

This is why the following code has a silent bug:

```go
// BAD
func findUser(found bool) Stringer {
    var u *User              // type: *User, value: nil
    if !found {
        return u             // returns (*User, nil) — the interface is NOT nil
    }
    return &User{Name: "Alice"}
}

result := findUser(false)
if result == nil {           // FALSE — the check misses the bug
    fmt.Println("not found")
}
result.String()              // PANIC — value word is nil
```

---

## The three guards — and when to use each

### Guard 1 — Fix it at the source (always prefer this)

The root cause is that a typed nil escaped into an interface.
The fix is: **never return a typed nil variable through an interface — return `nil` directly**.

```go
// BAD — typed nil leaks in
func findUser(found bool) Stringer {
    var u *User
    if !found {
        return u   // (*User, nil) — bug
    }
    return &User{Name: "Alice"}
}

// GOOD — return an untyped nil
func findUser(found bool) Stringer {
    if !found {
        return nil  // (nil, nil) — correct, == nil check works
    }
    return &User{Name: "Alice"}
}
```

The untyped `nil` has no type, so Go assigns `(nil, nil)` to the interface — a true nil interface.
Now the caller's `== nil` check works correctly.

**When to use:** Any time you write a function that returns an interface type.
This should be your default habit.

---

### Guard 2 — Type assertion (preferred when you know the concrete type)

If you receive an interface from code you trust but cannot change, and you know what
concrete type it wraps, assert to that type and check the pointer directly:

```go
func safeCall(s Stringer) {
    u, ok := s.(*User)   // step 1: pull out the concrete *User
    if !ok || u == nil { // step 2: check both "was it a *User?" and "is the pointer nil?"
        fmt.Println("skipping — nil or wrong type")
        return
    }
    fmt.Println(u.String()) // safe — u is a non-nil *User
}
```

What happens under the hood during `s.(*User)`:

```
s = (*User, nil)

type assertion checks: is type word == *itab(User)? → YES → ok = true
extracts value word → nil pointer → u = nil
```

So `ok` is `true` (the type matched) but `u` is `nil` (the pointer is nil).
That's why you must check **both** `ok` and `u == nil`.

**When to use:** When you know what concrete type to expect, e.g. inside a service
that creates the interface values itself.

---

### Guard 3 — `reflect.ValueOf` (last resort, for unknown types)

When you receive an `any` / `interface{}` from truly unknown external code and you
cannot know the concrete type at compile time, use reflection to inspect the value word directly:

```go
func isTrulyNil(i any) bool {
    if i == nil {
        return true  // fast path: (nil, nil)
    }
    v := reflect.ValueOf(i)          // unwrap the interface, get a reflect.Value
    switch v.Kind() {                 // what kind of thing is stored?
    case reflect.Ptr, reflect.Interface,
         reflect.Slice, reflect.Map,
         reflect.Chan, reflect.Func:
        return v.IsNil()             // check the value word directly
    }
    return false  // int, string, struct etc. can never be nil
}
```

Step by step for `isTrulyNil(s)` where `s = (*User, nil)`:

```
i == nil?                  → false (type word is *User, not nil)
reflect.ValueOf(i)         → a reflect.Value that wraps the *User nil pointer
v.Kind()                   → reflect.Ptr
v.IsNil()                  → true  ← looks directly at the value word
return true ✅
```

Step by step for `isTrulyNil(42)`:

```
i == nil?                  → false
reflect.ValueOf(i)         → a reflect.Value wrapping int 42
v.Kind()                   → reflect.Int  (not in the switch cases)
return false ✅             (an int can never be nil)
```

**When to use:** Frameworks, libraries, or generic utilities that receive `any` and
cannot assume anything about the concrete type. It is the most expensive option
(reflection is slower than a direct assertion), so avoid it in hot paths.

---

## Decision tree: which guard to use?

```
You have an interface value and need to know if it's safe to use
│
├─ Can you change the code that PRODUCES the interface value?
│   └─ YES → Fix at the source: return nil, not a typed nil variable.  (Guard 1)
│
└─ NO, you receive it from outside
    │
    ├─ Do you know the concrete type?
    │   └─ YES → Type assertion: u, ok := s.(*YourType); ok && u != nil  (Guard 2)
    │
    └─ NO, it's truly dynamic / any
        └─ reflect.ValueOf(i).IsNil()                                   (Guard 3)
```

---

## One-line summary

> An interface is nil only when **both** the type word and the value word are zero.
> Assigning a nil pointer fills the type word, making the interface non-nil.
> The only real fix is to never let a typed nil enter an interface — return `nil` directly.

