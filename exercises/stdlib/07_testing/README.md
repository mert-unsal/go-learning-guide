# 📦 Module 07 — testing: Go Testing Patterns & Techniques

> **Topics covered:** table-driven tests · subtests · t.Parallel · benchmarks · panic testing · test helpers · errors.As · function injection · t.Helper
>
> **Deep dive:** [Chapter 14 — Testing Internals](../../../learnings/14_testing_internals.md)

---

## 🗺️ Learning Path

```
1. Read: Chapter 14 — how go test works under the hood
2. Open exercises.go + exercises_test.go side by side
3. Study the TEST first — it shows the pattern
4. Implement the function to make the test pass
```

---

## 📚 Testing Patterns Demonstrated

| Pattern | Exercise | Key Insight |
|---------|----------|-------------|
| Table-driven tests + `t.Run` | Ex 1 | **The Go standard** for testing multiple inputs |
| Error path testing | Ex 2 | **Always test both happy and error paths** |
| Panic testing with `recover` | Ex 3 | **defer + recover to assert panics** |
| Benchmarking with `b.N` | Ex 4 | **go test -bench=. for performance measurement** |
| `t.Parallel()` subtests | Ex 5 | **Concurrent test execution** |
| Golden value tests | Ex 6 | **Multi-byte rune correctness** |
| `t.Helper()` | Ex 7 | **Custom assert functions** with correct line numbers |
| `errors.As` type assertion | Ex 8 | **Test custom error types** in the chain |
| Immutability verification | Ex 9 | **Ensure original data isn't modified** |
| Function injection (no mocks) | Ex 10 | **Inject sleep/time for testability** |
| Map-based test cases | Ex 11 | **Concise lookup-table testing** |
| Higher-order function testing | Ex 12 | **Pass any func(string) string** |

---

## ✏️ Exercises

| # | Function | Testing Pattern |
|---|----------|----------------|
| 1 | `AddEx(a, b)` | Table-driven + `t.Run` + `t.Parallel` |
| 2 | `DivideEx(a, b)` | Happy path + error case |
| 3 | `MaxEx(nums)` | Panic testing with `recover` |
| 4 | `ContainsEx(s, target)` | Benchmark + `b.N` |
| 5 | `FizzBuzzEx(n)` | Parallel subtests |
| 6 | `ReverseEx(s)` | Golden values (including Unicode) |
| 7 | `IsSortedEx(nums)` | `t.Helper()` assertion |
| 8 | `ParseKeyValue(s)` | `errors.As` + custom error type |
| 9 | `SortStringsEx(input)` | Verify no mutation of input |
| 10 | `Retry(max, sleep, fn)` | Injected dependencies |
| 11 | `HTTPStatusText(code)` | Map-based test cases |
| 12 | `Transform(input, fn)` | Function injection |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/stdlib/07_testing/

# Run benchmarks
go test -bench=. -benchmem ./exercises/stdlib/07_testing/
```

---

## ✅ Done? Next Step

```bash
go test -race -v ./exercises/stdlib/08_context/
```
