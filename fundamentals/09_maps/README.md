# 📦 Module 09 — Maps

> **Topics covered:** Map creation · Iteration · Existence checks · Frequency counting · Grouping · Anagram detection

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
| `make(map[K]V)` and map literals | `concepts.go` |
| Safe existence check `v, ok := m[key]` | Exercise 3 — `SafeGet` in 07 |
| Iterating with `for range` | Exercise 1 — `CharFrequency` |
| Grouping data with maps | Exercise 2 — `GroupByFirstChar` |
| Frequency counting | Exercise 3 — `TopTwoFrequent` |
| Using a map as a "seen" set | Exercise 5 — `FirstDuplicate` |
| Word counting | Exercise 6 — `WordCount` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `CharFrequency(s string) map[rune]int` | Count each character in a string |
| 2 | `GroupByFirstChar(words []string) map[byte][]string` | Group words by first letter |
| 3 | `TopTwoFrequent(nums []int) []int` | Find the two most frequent numbers |
| 4 | `IsAnagram(s, t string) bool` | Check if two strings are anagrams |
| 5 | `FirstDuplicate(nums []int) int` | Find the first number that appears twice |
| 6 | `WordCount(sentence string) map[string]int` | Count occurrences of each word |

---

## 🧪 Run Tests

> ⚠️ The `./fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./fundamentals/09_maps/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestCharFrequency
go test . -v -run TestGroupByFirstChar
go test . -v -run TestTopTwoFrequent
go test . -v -run TestIsAnagram
go test . -v -run TestFirstDuplicate
go test . -v -run TestWordCount
```

---

## 💡 Key Hints

<details>
<summary>Exercise 1 — CharFrequency hint</summary>

`for range` on a string gives `(index, rune)` pairs:
```go
freq := make(map[rune]int)
for _, ch := range s {
    freq[ch]++   // safe — zero value of int is 0
}
return freq
```
</details>

<details>
<summary>Exercise 4 — IsAnagram hint</summary>

Count chars in `s`, subtract for `t`, check everything is zero:
```go
counts := make(map[rune]int)
for _, ch := range s { counts[ch]++ }
for _, ch := range t { counts[ch]-- }
for _, v := range counts {
    if v != 0 { return false }
}
return true
```
</details>

<details>
<summary>Important: Map zero value</summary>

If you access a key that doesn't exist, Go returns the **zero value** (0 for int, "" for string). This means `m[key]++` works even for new keys — no need to check if the key exists first!
</details>

---

## ✅ Done? Next Step

```bash
go test ./fundamentals/10_goroutines/... -v
```

