# 📦 Module 03 — Functions

> **Topics covered:** Multiple return values · Variadic functions · Higher-order functions · Closures · Recursion · Memoization

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
| Multiple return values | `concepts.go` + Exercise 1 |
| Variadic functions (`...int`) | `concepts.go` + Exercise 2 |
| Functions as values (first-class) | `concepts.go` + Exercise 3 |
| Higher-order functions | Exercise 3 — `Apply` |
| Closures (capturing variables) | `concepts.go` + Exercise 4 |
| Recursion | Exercise 5 — `Fibonacci` |
| Memoization with a map | Exercise 5 — `FibonacciMemo` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `MinMax(nums []int) (min, max int)` | Return min and max of a slice using multiple return values |
| 2 | `Sum(nums ...int) int` | Variadic sum — `Sum(1,2,3)` → 6 |
| 3 | `Apply(nums []int, fn func(int) int) []int` | Apply a function to each element of a slice |
| 4 | `MakeAdder(n int) func(int) int` | Return a closure that adds `n` to its argument |
| 5 | `Fibonacci(n int) int` | Recursive Fibonacci |
| 5b | `FibonacciMemo(n int) int` | Memoized Fibonacci using a `map[int]int` |

---

## 🧪 Run Tests

### Run all tests for this module:
```bash
go test ./fundamentals/03_functions/... -v
```

### Run a single exercise test:
```bash
go test ./fundamentals/03_functions/... -v -run TestMinMax
go test ./fundamentals/03_functions/... -v -run TestSum
go test ./fundamentals/03_functions/... -v -run TestApply
go test ./fundamentals/03_functions/... -v -run TestMakeAdder
go test ./fundamentals/03_functions/... -v -run TestFibonacci
go test ./fundamentals/03_functions/... -v -run TestFibonacciMemo
```

---

## 💡 Key Hints

<details>
<summary>Exercise 3 — Apply hint</summary>

Create a new slice and apply `fn` to each element:
```go
result := make([]int, len(nums))
for i, v := range nums {
    result[i] = fn(v)
}
return result
```
</details>

<details>
<summary>Exercise 4 — MakeAdder closure hint</summary>

The returned function "captures" the variable `n`:
```go
func MakeAdder(n int) func(int) int {
    return func(x int) int {
        return x + n   // n is captured from outer scope
    }
}
```
</details>

<details>
<summary>Exercise 5b — Memoization hint</summary>

Check the map before computing recursively:
```go
if val, ok := memo[n]; ok {
    return val
}
result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
memo[n] = result
return result
```
</details>

---

## ✅ Done? Next Step

```bash
go test ./fundamentals/04_pointers/... -v
```

