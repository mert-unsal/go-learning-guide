# 📦 Module 06 — Interfaces

> **Topics covered:** Interface definition · Implicit implementation · Polymorphism · Type assertions · Type switches · `io.Writer` pattern

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
| Defining and implementing interfaces | Exercise 1 — `ExStringer` |
| Passing an interface to a function | Exercise 1 — `PrintAll` |
| The `io.Writer` pattern (custom writers) | Exercise 2 — `ExWriter` |
| Type assertions `x.(T)` | `concepts.go` |
| Type switches `switch v := x.(type)` | Exercise 3 — `Describe` |
| Empty interface `interface{}` / `any` | `concepts.go` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1a | `(b ExBook) String() string` | Return `"Title" by Author` |
| 1b | `(m ExMovie) String() string` | Return `Title (Year)` |
| 1c | `PrintAll(items []ExStringer)` | Call `item.String()` on each |
| 2a | `(bw *ExBufferWriter) Write(data string) error` | Append data to `Buffer` |
| 2b | `WriteAll(w ExWriter, items []string) error` | Call `w.Write` for each item |
| 3 | `Describe(i interface{}) string` | Type switch on int, string, bool |

---

## 🧪 Run Tests

> ⚠️ The `./fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./fundamentals/06_interfaces/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestStringer
go test . -v -run TestWriter
go test . -v -run TestDescribe
```

---

## 💡 Key Hints

<details>
<summary>Exercise 3 — Type switch hint</summary>

```go
func Describe(i interface{}) string {
    switch v := i.(type) {
    case int:
        return fmt.Sprintf("int: %d", v)
    case string:
        return fmt.Sprintf("string: %s", v)
    case bool:
        return fmt.Sprintf("bool: %t", v)
    default:
        return "unknown"
    }
}
```
</details>

<details>
<summary>Go interfaces are implicit — key insight</summary>

Unlike Java/C#, you never write `implements ExStringer`. If your type has the right methods, it **automatically** satisfies the interface. This enables powerful decoupling — packages don't need to import each other to share interfaces.
</details>

---

## ✅ Done? Next Step

```bash
go test ./fundamentals/07_error_handling/... -v
```

