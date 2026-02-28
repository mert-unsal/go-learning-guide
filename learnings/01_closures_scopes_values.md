# Deep Dive: Scopes, Closures, and How Go Passes Things

---

## Question 1: How can the compiler allow `i := i` when `i` already exists?

### The answer: **Scopes (Blocks)**

Go uses **block scoping**. Every `{...}` creates a new scope. A variable declared inside a block only lives in that block and its children.

```go
i := 5          // outer scope (the for loop's "init" scope)

{               // new inner block (the loop body)
    i := i      // NEW variable named i, initialized from outer i
    // from here, "i" refers to the INNER variable
    fmt.Println(i) // 5
}

fmt.Println(i) // still 5 — outer i is unchanged
```

The compiler is perfectly fine with this because:
- The **outer `i`** lives in the loop's outer scope
- The **inner `i`** lives in the loop body's scope
- They are **two separate memory locations** with the same name
- The inner one **shadows** (hides) the outer one within its block

Think of it like two rooms:
```
[Outer Room: i = 0]
    └── [Inner Room: i = 0]  ← a completely different "i" box
```

The compiler always knows which `i` you mean based on which room (scope) you're in.

---

## Question 2: Does Go pass a reference in the for loop?

### Short answer: **No. Go always passes VALUES. But closures capture VARIABLES.**

This is the most important distinction to understand.

---

### What "capturing a variable" means

A **closure** is a function that "closes over" variables from its surrounding scope. It does NOT copy the value — it keeps a **live link** to the original variable itself.

```go
x := 10
f := func() {
    fmt.Println(x)  // f "closes over" x — it holds a link to x
}
x = 99
f() // prints 99 ❌ — f sees the CURRENT value of x, not 10
```

The function `f` doesn't hold the value `10`. It holds a pointer to the variable `x`. Whatever `x` is when you CALL `f`, that's what it prints.

---

### Why all closures print `3` in the loop

```go
for i := 0; i < 3; i++ {
    funcs[i] = func() {
        fmt.Println(i)  // all 3 closures point to THE SAME i
    }
}
```

There is only **one `i` variable** created for the entire loop. All 3 closures point to this single variable.

Memory picture:

```
                    ┌─────────────┐
funcs[0] ──────────►│             │
funcs[1] ──────────►│  prints: i  │──── all pointing to ──►  [ i = 3 ]
funcs[2] ──────────►│             │
                    └─────────────┘
```

When the loop ends, `i` becomes `3`. So when you call any of the functions, they all read `i` and see `3`.

---

### With `i := i` — each closure gets its own variable

```go
for i := 0; i < 3; i++ {
    i := i  // creates a NEW variable per iteration
    funcs[i] = func() {
        fmt.Println(i)  // each closure points to its OWN i
    }
}
```

Memory picture:

```
funcs[0] ──► closure ──► [ i = 0 ]   (its own box, frozen at 0)
funcs[1] ──► closure ──► [ i = 1 ]   (its own box, frozen at 1)
funcs[2] ──► closure ──► [ i = 2 ]   (its own box, frozen at 2)
```

3 separate variables. 3 separate closures. Each pointing to a different box.

---

## Question 3: "Arguments are always copied" — what does that mean?

### Go is a **pass-by-value** language

When you call a function with arguments, Go **copies** the value into a new variable for the function to use. The original is untouched.

```go
func double(n int) {
    n = n * 2
    fmt.Println(n) // 20
}

x := 10
double(x)
fmt.Println(x) // still 10 ✅ — x was COPIED, not handed over
```

`double` gets its own copy of `x`. Whatever it does to `n` doesn't affect `x`.

---

### Why this fixes the closure problem

```go
for i := 0; i < 3; i++ {
    go func(n int) {   // n is a COPY of i at this moment
        fmt.Println(n)
    }(i)               // i is passed here — copied into n
}
```

- At iteration 0: `i=0` is **copied** into `n`. Closure captures `n` which is its own private copy, value `0`.
- At iteration 1: `i=1` is **copied** into `n`. New `n`, value `1`.
- At iteration 2: `i=2` is **copied** into `n`. New `n`, value `2`.

Each function invocation gets its own parameter `n` — a completely independent copy.

---

### The contrast: closures vs arguments

| | Closure (captures variable) | Argument (passed as copy) |
|---|---|---|
| What it holds | A **live link** to the original variable | A **snapshot copy** of the value |
| What happens if original changes | Sees the **new value** | Doesn't care — has its own copy |
| Like... | A camera pointing at a scoreboard | A photograph of the scoreboard |

---

## Summary

```
SCOPE:     { } creates a new scope. Same name in inner scope = new variable. Shadowing.

CLOSURE:   captures a VARIABLE (live link), not a value.
           → if the variable changes, the closure sees the change

ARGUMENT:  always a COPY of the value at the moment of the call.
           → changes to the original don't affect the function's copy
```

The loop bug happens because:
1. There's only **one `i` variable** in the loop
2. Closures hold a **live link** to it (not a copy)
3. By call-time, `i` has been incremented to its final value

The fix (`i := i` or passing as argument) works because:
1. It creates a **new variable** per iteration
2. That new variable is **never modified** after creation
3. So the closure always sees the correct snapshot value

