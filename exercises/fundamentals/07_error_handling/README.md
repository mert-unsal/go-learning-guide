# 📦 Module 07 — Error Handling

> **Topics covered:** `error` interface · Custom error types · Sentinel errors · `errors.Is` / `errors.As` · Error wrapping with `%w`

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
| Returning `error` as second value | Exercise 1 — `Divide` |
| Custom error type with `Error() string` | Exercise 2 — `ValidationError` |
| Sentinel errors (`errors.New`) | Exercise 4 — `ErrUserNotFound` |
| Error wrapping with `fmt.Errorf("%w")` | Exercise 5 — `WrapError` |
| `errors.Is` and `errors.As` | `concepts.go` |
| Safe map access with error | Exercise 3 — `SafeGet` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `Divide(a, b float64) (float64, error)` | Return error when `b == 0` |
| 2 | `Validate(name string) error` | Return `*ValidationError` if name is empty |
| 3 | `SafeGet(m map[string]int, key string) (int, error)` | Return error if key not in map |
| 4 | `FindUser(id int) (string, error)` | Return sentinel errors for id≤0 and id==999 |
| 5 | `WrapError(err error, context string) error` | Wrap with `fmt.Errorf("%s: %w", context, err)` |

---

## 🧪 Run Tests

> ⚠️ The `./exercises/fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./exercises/fundamentals/07_error_handling/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestDivide
go test . -v -run TestValidate
go test . -v -run TestSafeGet
go test . -v -run TestFindUser
go test . -v -run TestWrapError
```

---

## 💡 Key Hints

<details>
<summary>Go error handling philosophy</summary>

Go does NOT use exceptions. Instead, functions return errors as a second value:
```go
result, err := SomeFunction()
if err != nil {
    // handle it — don't ignore!
    return fmt.Errorf("context: %w", err)
}
```
Always handle errors explicitly. Ignoring them is a common bug.
</details>

<details>
<summary>Exercise 2 — Custom error type hint</summary>

The `ValidationError` struct is already defined. Just implement `Validate`:
```go
func Validate(name string) error {
    if name == "" {
        return &ValidationError{Field: "name", Message: "cannot be empty"}
    }
    return nil
}
```
</details>

<details>
<summary>Exercise 5 — Why wrap errors with %w?</summary>

`%w` lets callers use `errors.Is(err, ErrOriginal)` to unwrap and check the original error. Without `%w`, the original error is hidden inside a string.
</details>

---

## ✅ Done? Next Step

```bash
go test ./exercises/fundamentals/08_arrays_slices/... -v
```

