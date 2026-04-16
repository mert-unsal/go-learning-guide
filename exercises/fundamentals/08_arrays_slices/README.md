# 📦 Module 08 — Arrays & Slices

> **Topics covered:** Arrays vs slices · `make` · `append` · Two-pointer technique · In-place operations · 2D slices

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
| Array vs Slice difference | `concepts.go` |
| `make([]int, len, cap)` | Exercise 3 — `Make2D` |
| `append` and slice growth | `concepts.go` |
| In-place modification | Exercise 1 — `ReverseSlice` |
| Two-pointer technique | Exercise 1, 2 |
| Write-pointer pattern | Exercise 2 — `RemoveDuplicates` |
| Three-reversal trick | Exercise 4 — `RotateLeft` |
| Higher-order slice functions | Exercise 5 — `Filter` |
| Two-pointer merge | Exercise 6 — `MergeSorted` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `ReverseSlice(s []int)` | Reverse in-place using two pointers |
| 2 | `RemoveDuplicates(s []int) []int` | Remove duplicates from sorted slice in-place |
| 3 | `Make2D(rows, cols int) [][]int` | Allocate a rows×cols 2D slice |
| 4 | `RotateLeft(s []int, k int)` | Rotate left by k using the 3-reversal trick |
| 5 | `Filter(s []int, fn func(int) bool) []int` | Return elements that satisfy fn |
| 6 | `MergeSorted(a, b []int) []int` | Merge two sorted slices into one |

---

## 🧪 Run Tests

> ⚠️ The `./exercises/fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./exercises/fundamentals/08_arrays_slices/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestReverseSlice
go test . -v -run TestRemoveDuplicates
go test . -v -run TestMake2D
go test . -v -run TestRotateLeft
go test . -v -run TestFilter
go test . -v -run TestMergeSorted
```

---

## 💡 Key Hints

<details>
<summary>Exercise 1 — ReverseSlice two-pointer hint</summary>

```go
left, right := 0, len(s)-1
for left < right {
    s[left], s[right] = s[right], s[left]
    left++
    right--
}
```
</details>

<details>
<summary>Exercise 2 — RemoveDuplicates write-pointer hint</summary>

Use a write pointer `w` that only advances when a new value is seen:
```go
w := 1
for r := 1; r < len(s); r++ {
    if s[r] != s[r-1] {
        s[w] = s[r]
        w++
    }
}
return s[:w]
```
</details>

<details>
<summary>Exercise 4 — RotateLeft three-reversal hint</summary>

For `[1,2,3,4,5]` with k=2:
1. Reverse all → `[5,4,3,2,1]`
2. Reverse first `n-k` → `[3,4,5,2,1]`
3. Reverse last `k` → `[3,4,5,1,2]`

Remember: normalize `k = k % len(s)` first!
</details>

---

## ✅ Done? Next Step

```bash
go test ./exercises/fundamentals/09_maps/... -v
```

---

## 📖 Companion Chapter

For the deep-dive theory behind these exercises, read:

- [01 — Slices: The Three-Word Header](../../../learnings/01_slices_three_word_header.md) — slice header struct, growslice, backing array sharing, capacity traps

