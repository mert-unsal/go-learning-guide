# 📦 Module 09 — Maps

> **Topics covered:** hash tables · nil map safety · comma-ok pattern · map-as-set · iteration randomness · merge patterns
>
> **Deep dive:** [Chapter 02 — Maps: Buckets, Growth & the Never-Shrink Truth](../../../learnings/02_maps_buckets_and_growth.md)

---

## 🗺️ Learning Path

```
1. Read learnings/02_maps_buckets_and_growth.md   ← How runtime.hmap works under the hood
2. Run cmd/concepts/maps/*                        ← See map behavior in action
3. Open exercises.go                               ← Implement the 12 exercises
4. Run go test -race -v ./...                      ← Make them all pass
```

---

## 📚 What You Will Learn

| Concept | Exercise | Under the Hood |
|---------|----------|---------------|
| Frequency counting with maps | Ex 1 | Zero value of int is 0 → `m[key]++` works on new keys |
| Map of slices (grouping) | Ex 2 | `append(nil, x)` works — nil slice is valid |
| Two-pass frequency scan | Ex 3 | Build frequency map, then scan for top-N |
| Anagram detection | Ex 4 | Count up / count down pattern |
| Map as "seen" set | Ex 5 | `map[T]bool` vs `map[T]struct{}` tradeoff |
| Word frequency | Ex 6 | `strings.Fields` vs `strings.Split` |
| **Nil map safety** | Ex 7 | Reading nil map = safe. Writing = **panic** |
| **Map inversion** | Ex 8 | Multiple keys → same value → `map[V][]K` |
| **Merge with collision resolution** | Ex 9 | No built-in merge — iterate and decide |
| **Set difference** | Ex 10 | `map[T]bool` as set, `if b[key]` for membership |
| **Unique values + sort** | Ex 11 | Map iteration order is **random** — sort for determinism |
| **Map equality** | Ex 12 | Maps can't use `==` — must check manually |

---

## ✏️ Exercises

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `CharFrequency(s)` | Count each character in a string |
| 2 | `GroupByFirstChar(words)` | Group words by first letter |
| 3 | `TopTwoFrequent(nums)` | Find the two most frequent numbers |
| 4 | `IsAnagram(s, t)` | Check if two strings are anagrams |
| 5 | `FirstDuplicate(nums)` | First number that appears twice |
| 6 | `WordCount(sentence)` | Count word occurrences |
| 7 | `NilMapRead(m, key)` | **Comma-ok on nil map** |
| 8 | `InvertMap(m)` | **Swap keys and values** |
| 9 | `MergeMaps(a, b, resolve)` | **Merge with collision function** |
| 10 | `SetDifference(a, b)` | **Keys in a but not in b** |
| 11 | `UniqueValues(m)` | **Deduplicate and sort values** |
| 12 | `MapEqual(a, b)` | **Manual equality (no == for maps)** |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/fundamentals/09_maps/
```

---

## ✅ Done? Next Step

```bash
go test -race -v ./exercises/fundamentals/10_goroutines/
```

