# The Classic Closure Loop Gotcha in Go

---

## The Problem

```go
funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    funcs[i] = func() {
        fmt.Println(i)
    }
}
funcs[0]() // prints 3 ❌ (not 0!)
funcs[1]() // prints 3 ❌ (not 1!)
funcs[2]() // prints 3 ❌ (not 2!)
```

**Why?** All 3 functions are closures that *capture the same variable `i`* — not its value, but a **reference** to the variable. By the time you call them, the loop is done and `i` is already `3`.

Think of it like 3 people all looking at the same scoreboard. When the game ends, everyone sees the final score: 3.

---

## The Fix: `i := i`

```go
for i := 0; i < 3; i++ {
    i := i  // ← THIS LINE
    funcs[i] = func() {
        fmt.Println(i)
    }
}
funcs[0]() // prints 0 ✅
funcs[1]() // prints 1 ✅
funcs[2]() // prints 2 ✅
```

**What does `i := i` do?**

- The `:=` creates a **brand new variable** called `i` inside the loop body.
- This new `i` is a **copy** of the current loop `i` at that iteration.
- Each closure now captures its **own private copy** of `i`.

It's called **variable shadowing** — the inner `i` shadows (hides) the outer loop `i` within that block.

---

## Step by Step

| Iteration | Loop `i` | Inner `i` (new copy) | Closure captures |
|-----------|----------|----------------------|-----------------|
| 1st       | 0        | 0 (new var)          | its own `i = 0` |
| 2nd       | 1        | 1 (new var)          | its own `i = 1` |
| 3rd       | 2        | 2 (new var)          | its own `i = 2` |

Each closure gets its own separate box to store its value.

---

## Alternative Fix (also common)

Pass `i` as a function argument — arguments are always copied:

```go
for i := 0; i < 3; i++ {
    func(i int) {          // i is now a parameter (a copy)
        funcs[i] = func() {
            fmt.Println(i)
        }
    }(i)                   // immediately invoke with current i
}
```

---

## Key Takeaway

> Closures capture **variables** (references), not **values**.  
> If the variable changes after the closure is created, the closure sees the new value.  
> `i := i` forces a new variable per iteration so each closure has its own independent copy.

