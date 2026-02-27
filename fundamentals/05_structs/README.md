# 📦 Module 05 — Structs

> **Topics covered:** Struct definition · Methods · Interfaces · Embedding · Stack data structure

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
| Struct definition and initialization | `concepts.go` |
| Methods on structs | Exercise 1 — `ExRectangle` |
| Implementing interfaces | Exercise 2 — `ExShape` |
| Polymorphism via interfaces | Exercise 2 — `TotalArea` |
| Building a data structure with a struct | Exercise 3 — `ExStack` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1a | `(r ExRectangle) Area() float64` | `Width * Height` |
| 1b | `(r ExRectangle) Perimeter() float64` | `2 * (Width + Height)` |
| 2 | `TotalArea(shapes []ExShape) float64` | Sum the area of all shapes |
| 3a | `(s *ExStack) Push(val int)` | Append to the internal slice |
| 3b | `(s *ExStack) Pop() (int, bool)` | Remove and return the top element |
| 3c | `(s *ExStack) IsEmpty() bool` | Return `len(items) == 0` |

---

## 🧪 Run Tests

> ⚠️ The `./fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./fundamentals/05_structs/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestRectangle
go test . -v -run TestTotalArea
go test . -v -run TestStack
```

---

## 💡 Key Hints

<details>
<summary>Exercise 3 — Stack with a slice hint</summary>

Add a `items []int` field to the struct first:
```go
type ExStack struct {
    items []int
}
```

Then:
- `Push`: `s.items = append(s.items, val)`
- `Pop`: grab the last element, shrink the slice — check if empty first!
```go
func (s *ExStack) Pop() (int, bool) {
    if s.IsEmpty() { return 0, false }
    top := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return top, true
}
```
</details>

<details>
<summary>ExRectangle satisfies ExShape automatically</summary>

In Go, interfaces are **implicit** — if `ExRectangle` has both `Area()` and `Perimeter()`, it automatically satisfies `ExShape`. No `implements` keyword needed.
</details>

---

## ✅ Done? Next Step

```bash
go test ./fundamentals/06_interfaces/... -v
```

