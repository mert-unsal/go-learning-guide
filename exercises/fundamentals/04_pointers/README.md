# 📦 Module 04 — Pointers

> **Topics covered:** Pointer basics · Dereferencing · Pointer receivers · Value receivers · Constructor pattern

---

## 🗺️ Learning Path

```
1. Read concepts.go        ← Theory with runnable examples
2. Open exercises.go       ← Implement the TODO functions yourself
3. Run the tests below     ← Instant feedback on your code
4. Stuck? Open solutions.go ← Only after you have tried!
```

---

## 📚 What You Will Learn

| Concept | Where |
|---------|-------|
| `&` (address-of) and `*` (dereference) | `concepts.go` |
| Passing pointers to functions | Exercise 1 — `Increment` |
| Swapping via pointers | Exercise 2 — `SwapPointers` |
| Pointer receiver vs value receiver | Exercise 3 — `ScoreBoard` |
| Constructor pattern (`NewXxx`) | Exercise 4 — `NewPlayer` |
| Modifying values through pointers | Exercise 5 — `DoubleValue` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `Increment(n *int)` | Increment the integer at pointer address: `*n++` |
| 2 | `SwapPointers(a, b *int)` | Swap two values via pointers: `*a, *b = *b, *a` |
| 3a | `(s *ScoreBoard) AddPoints(points int)` | Pointer receiver — mutates the struct |
| 3b | `(s ScoreBoard) CurrentScore() int` | Value receiver — just reads the struct |
| 4 | `NewPlayer(name string, level int) *Player` | Return `&Player{Name: name, Level: level}` |
| 5 | `DoubleValue(x *int)` | Double the value at pointer: `*x = *x * 2` |

---

## 🧪 Run Tests

> ⚠️ The `./exercises/fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./exercises/fundamentals/04_pointers/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestIncrement
go test . -v -run TestSwapPointers
go test . -v -run TestScoreBoard
go test . -v -run TestNewPlayer
go test . -v -run TestDoubleValue
```

---

## 💡 Key Hints

<details>
<summary>Pointer receiver vs value receiver — when to use which?</summary>

| Use **pointer receiver** `*T` | Use **value receiver** `T` |
|-------------------------------|---------------------------|
| You need to **mutate** the struct | You only **read** the struct |
| Struct is large (avoid copying) | Struct is small & immutable |
| `AddPoints` ← changes Score | `CurrentScore` ← just returns Score |

</details>

<details>
<summary>Exercise 4 — Constructor pattern hint</summary>

In Go, there's no `new` keyword like Java/C#. The convention is:
```go
func NewPlayer(name string, level int) *Player {
    return &Player{Name: name, Level: level}
}
```
The `&` in front of the struct literal creates it on the heap and returns a pointer.
</details>

---

## ✅ Done? Next Step

```bash
go test ./exercises/fundamentals/05_structs/... -v
```

