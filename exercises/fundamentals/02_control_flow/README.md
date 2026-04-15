# 📦 Module 02 — Control Flow

> **Topics covered:** defer internals · range rewrites · switch dispatch · labeled breaks · Go 1.22 changes
>
> **Deep dive:** [Chapter 22 — Control Flow Under the Hood](../../../learnings/22_control_flow_under_the_hood.md)

---

## 🗺️ Learning Path

```
1. Read learnings/22_control_flow_under_the_hood.md    ← How the compiler transforms your code
2. Run cmd/concepts/control-flow/*                     ← See it in action
3. Open exercises.go                                   ← Implement the 12 exercises
4. Run go test -race -v ./...                          ← Make them all pass
```

---

## 📚 What You Will Learn

| Concept | Exercise | Under the Hood |
|---------|----------|---------------|
| `switch` (no fallthrough by default) | Ex 1 | Jump table vs linear scan |
| `for` loop (the only loop in Go) | Ex 2 | Bounds check elimination |
| `for range` over strings | Ex 3, 9 | Rune decoding, byte offset vs rune index |
| Early return & `break` inside loops | Ex 4 | — |
| `defer` — LIFO execution order | Ex 5 | `_defer` linked list on `runtime.g` |
| `defer` — named return modification | Ex 6 | Return value set → defer runs → caller gets it |
| `defer` — argument evaluation timing | Ex 7 | Args evaluated at defer-time, not execution-time |
| `range` — value copy semantics | Ex 8 | Compiler rewrites `v := slice[i]` (copy) |
| Labeled break for nested loops | Ex 10 | `break outer` exits the named loop |
| Type switch | Ex 11 | Pointer comparison on `_type` descriptor |
| `range` over integer (Go 1.22+) | Ex 12 | Rewrites to `for i := 0; i < n; i++` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `FizzBuzzSwitch(n)` | Classic FizzBuzz using `switch` |
| 2 | `SumTo(n)` | Sum 1..n with a `for` loop |
| 3 | `CountVowels(s)` | Count vowels using `for range` |
| 4 | `IsPrime(n)` | Primality check with early return |
| 5 | `DeferOrder()` | Predict LIFO defer execution order |
| 6 | `DeferModifyReturn(n)` | **Named return + defer modification** |
| 7 | `DeferArgCapture()` | **Defer argument evaluation timing** |
| 8 | `DoubleScores(players)` | **Range value copy trap** — modify originals |
| 9 | `RuneValues(s)` | **Range over string** — runes not bytes |
| 10 | `FindInMatrix(matrix, target)` | **Labeled break** in nested loops |
| 11 | `TypeDescribe(v)` | **Type switch** on interface values |
| 12 | `Squares(n)` | **Range over int** (Go 1.22+) |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/fundamentals/02_control_flow/
```

---

## ✅ Done? Next Step

```bash
go test -race -v ./exercises/fundamentals/03_functions/
```

