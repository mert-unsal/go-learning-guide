# Go's Automatic Pointer Magic: Auto-Address & Auto-Dereference

Go's compiler does **two silent tricks** to make your code less verbose. Understanding exactly when and why they happen is the key.

---

## The Two Tricks

| Name | What Go does silently | Example |
|---|---|---|
| **Auto-address** | `c.Method()` → `(&c).Method()` | When method needs `*T` but you have `T` |
| **Auto-deref** | `p.Field` → `(*p).Field` | When accessing fields on a `*T` |

---

## Trick 1: Auto-Dereference (field & method access on pointers)

### Rule: When you access a field or method through a pointer, Go silently inserts `*`

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

### When does auto-deref happen?

✅ Always — whenever you use `.` (dot notation) on a pointer type.

```go
p.Name    // ✅ auto-deref → (*p).Name
p.Age++   // ✅ auto-deref → (*p).Age++
p.Method() // ✅ auto-deref if Method has value receiver
```

### When must you write `*` explicitly?

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

## Trick 2: Auto-Address (calling pointer-receiver methods on values)

### Rule: When a method requires `*T` but you have `T`, Go silently inserts `&` — BUT only if the variable is addressable

```go
type Counter struct{ count int }

func (c *Counter) Increment() { c.count++ }  // pointer receiver

c := Counter{count: 0}  // c is a plain Counter (value, not pointer)

c.Increment()     // You write this
(&c).Increment()  // Go compiles it as this — c is addressable ✅
```

### What is "addressable"?

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

### The fix for non-addressable cases:

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

## The Full Decision Table

```
You have T,  method wants T  → direct call, nothing special
You have T,  method wants *T → Go auto-addresses: (&T).Method()  [if addressable]
You have *T, method wants *T → direct call, nothing special
You have *T, method wants T  → Go auto-derefs: (*T).Method()
```

---

## When to Write It Explicitly vs Leave It Implicit

### Leave it implicit (let Go handle it):

```go
// Field access — always let Go handle it
p.Name = "Alice"       // NOT (*p).Name = "Alice"
p.Age++                // NOT (*p).Age++

// Method calls on variables — always let Go handle it
c.Increment()          // NOT (&c).Increment()
p.SomeMethod()         // NOT (*p).SomeMethod()
```

Writing `(*p).Age` or `(&c).Increment()` is valid Go, but it's **noisy and unusual** — no Go developer writes it that way in practice.

### Write it explicitly when:

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

## One Slide Summary

```
AUTO-DEREF  (p.Field, p.Method)
  → Go inserts * automatically
  → ALWAYS works on pointer types with dot notation
  → Write * explicitly only when grabbing/replacing the WHOLE value

AUTO-ADDRESS  (c.PointerMethod())
  → Go inserts & automatically
  → ONLY works when variable is ADDRESSABLE (has a memory slot)
  → Fails on: temporaries, function returns, map values
  → Write & explicitly when passing to functions or creating pointers

GOLDEN RULE:
  Use dot notation normally (p.Name, c.Increment())
  Use * only when you mean "give me everything at this address"
  Use & only when you mean "give me the address of this thing"
```

