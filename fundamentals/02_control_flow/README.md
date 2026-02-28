# 📦 Module 02 — Control Flow

> **Topics covered:** if/else · switch · for loops · range · defer · break/continue

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
| `if / else if / else` | `concepts.go` |
| `switch` (no fallthrough by default) | `concepts.go` + Exercise 1 |
| `for` loop (the only loop in Go!) | `concepts.go` + Exercise 2 |
| `for range` over strings/slices | `concepts.go` + Exercise 3 |
| Early return & `break` inside loops | Exercise 4 |
| `defer` — LIFO execution order | `concepts.go` + Exercise 5 |
| `defer` — arguments evaluated at defer time | `concepts.go` → `DemonstrateDefereArgumentEvaluation()` |
| `defer` — can modify named return values | `concepts.go` → `deferWithNamedReturn()` |
| `defer` — loop gotcha & the fix | `concepts.go` → `deferInLoop()` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `FizzBuzzSwitch(n int) string` | Classic FizzBuzz — must use `switch`, not `if/else` |
| 2 | `SumTo(n int) int` | Sum integers 1..n using a `for` loop |
| 3 | `CountVowels(s string) int` | Count vowels a,e,i,o,u using `for range` (case-insensitive) |
| 4 | `IsPrime(n int) bool` | Return true if n is prime — use early `return` inside loop |
| 5 | `DeferOrder() []string` | Return `["third","second","first"]` — demonstrates defer LIFO |

---

## 🧪 Run Tests

> ⚠️ **Important:** The `./fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./fundamentals/02_control_flow/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from project root):
```bash
go test ./fundamentals/02_control_flow/... -v -run TestFizzBuzzSwitch
go test ./fundamentals/02_control_flow/... -v -run TestSumTo
go test ./fundamentals/02_control_flow/... -v -run TestCountVowels
go test ./fundamentals/02_control_flow/... -v -run TestIsPrime
go test ./fundamentals/02_control_flow/... -v -run TestDeferOrder
```

### From inside this folder:
```bash
go test . -v -run TestFizzBuzzSwitch
```

---

## 💡 Key Hints

<details>
<summary>Exercise 1 — FizzBuzz with switch hint</summary>

Go's `switch` can take an expression:
```go
switch {
case n%15 == 0:
    return "FizzBuzz"
case n%3 == 0:
    return "Fizz"
// ...
}
```
</details>

<details>
<summary>Exercise 3 — CountVowels hint</summary>

`for range` over a string gives you runes. Lowercase it first:
```go
for _, ch := range strings.ToLower(s) {
    // ch is a rune
}
```
</details>

<details>
<summary>Exercise 4 — IsPrime hint</summary>

Check divisors from 2 up to √n. If any divide evenly, it's not prime:
```go
for i := 2; i*i <= n; i++ {
    if n%i == 0 { return false }
}
```
</details>

<details>
<summary>Exercise 5 — Defer LIFO hint</summary>

Defers execute in **reverse order** (last-in, first-out). The answer is simply:
```go
return []string{"third", "second", "first"}
```

📌 **Bonus — Interview traps to know (see `concepts.go`):**

**Trap 1 — Arguments are captured at defer time:**
```go
x := 10
defer fmt.Println("x =", x) // captures x=10 RIGHT NOW
x = 99
// Output: "x = 10", not 99 !
```

**Trap 2 — Defer can modify named return values:**
```go
func foo() (result string) {
    defer func() { result = "changed" }()
    return "original" // gets overwritten — final return is "changed"
}
```

**Trap 3 — Never defer inside a loop for resources:**
```go
for i := range files {
    defer file.Close() // BAD: all close at end of function, not each iteration
}
// FIX: wrap in an anonymous function:
for i := range files {
    func() {
        defer file.Close() // closes after each anonymous call
    }()
}
```
</details>

---

## ✅ Done? Next Step

```bash
go test ./fundamentals/03_functions/... -v
```

