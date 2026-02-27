# 📦 Module 01 — Basics

> **Topics covered:** Variables · Constants · Types · iota · Multiple return values · Unicode strings

---

## 🗺️ Learning Path

Follow this order — do **not** jump ahead:

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
| `var` vs `:=` declarations | `concepts.go` — DemonstrateVariables |
| Zero values | `concepts.go` — DemonstrateVariables |
| `const` and `iota` | `concepts.go` — Weekday / Permission |
| Basic types (`int`, `float64`, `string`, `bool`) | `concepts.go` — DemonstrateTypes |
| Type conversions | `concepts.go` — DemonstrateTypes |
| Multiple return values | `exercises.go` — SwapInts |
| Unicode vs byte length | `exercises.go` — CharacterCount |
| Custom types with `iota` | `exercises.go` — Direction |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `CelsiusToFahrenheit(c float64) float64` | Formula: `F = C * 9/5 + 32` |
| 2 | `SwapInts(a, b int) (int, int)` | Return b, a using multiple return values |
| 3 | `CharacterCount(s string) int` | Count **characters** not bytes — use `[]rune` |
| 4 | `MinMax(a, b, c int) (min, max int)` | Return the smallest and largest of 3 ints |
| 5 | `DirectionName(d Direction) string` | Use a `switch` to return the direction name |

---

## 🧪 Run Tests

### Run all tests for this module:
```bash
go test ./fundamentals/01_basics/... -v
```

### Run a single exercise test:
```bash
go test ./fundamentals/01_basics/... -v -run TestCelsiusToFahrenheit
go test ./fundamentals/01_basics/... -v -run TestSwapInts
go test ./fundamentals/01_basics/... -v -run TestCharacterCount
go test ./fundamentals/01_basics/... -v -run TestMinMax
go test ./fundamentals/01_basics/... -v -run TestDirectionName
```

### Expected output when all pass:
```
=== RUN   TestCelsiusToFahrenheit
    ✅ CelsiusToFahrenheit(0) = 32
    ✅ CelsiusToFahrenheit(100) = 212
    ✅ CelsiusToFahrenheit(-40) = -40
    ✅ CelsiusToFahrenheit(37) = 98.6
--- PASS: TestCelsiusToFahrenheit (0.00s)
=== RUN   TestSwapInts
    ✅ SwapInts(3,7) = (7,3)
--- PASS: TestSwapInts (0.00s)
=== RUN   TestCharacterCount
    ✅ CharacterCount("hello") = 5
    ✅ CharacterCount("世界") = 2
    ✅ CharacterCount("Hello世界") = 7
    ✅ CharacterCount("") = 0
--- PASS: TestCharacterCount (0.00s)
=== RUN   TestMinMax
    ✅ MinMax(3,1,2) = (1,3)
    ✅ MinMax(5,5,5) = (5,5)
    ✅ MinMax(-1,-5,0) = (-5,0)
--- PASS: TestMinMax (0.00s)
=== RUN   TestDirectionName
    ✅ DirectionName(0) = "North"
    ✅ DirectionName(1) = "East"
    ✅ DirectionName(2) = "South"
    ✅ DirectionName(3) = "West"
--- PASS: TestDirectionName (0.00s)
ok      gointerviewprep/fundamentals/01_basics
```

---

## 💡 Key Hints

<details>
<summary>Exercise 3 — CharacterCount hint</summary>

`len(s)` returns the number of **bytes**, not characters.  
A Chinese character like `世` takes **3 bytes** but is only **1 character**.  
Solution: convert the string to a rune slice first.

```go
return len([]rune(s))
```
</details>

<details>
<summary>Exercise 4 — MinMax hint</summary>

Use `if` comparisons. Start by assuming `a` is both min and max, then compare with `b` and `c`.

```go
min, max = a, a
if b < min { min = b }
// ... and so on
```
</details>

<details>
<summary>Exercise 5 — DirectionName hint</summary>

Use a `switch` statement on the Direction type:

```go
switch d {
case North:
    return "North"
// ...
}
```
</details>

---

## ✅ Done? Next Step

Once all 5 tests pass, move to the next module:

```bash
cd ../02_control_flow
```

Or run from the project root:
```bash
go test ./fundamentals/02_control_flow/... -v
```

