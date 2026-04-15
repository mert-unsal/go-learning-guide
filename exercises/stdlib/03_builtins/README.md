# 📦 Module 03 — builtins: Built-in Functions & Patterns

> **Topics covered:** deep copy · recover · flatten · unique · chunking
>
> **Deep dive:** [Chapter 01 — Slices](../../../learnings/01_slices_three_word_header.md), [Chapter 02 — Maps](../../../learnings/02_maps_buckets_and_growth.md)

---

## ✏️ Exercises

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `DeepCopySlice(src)` | Copy slice without sharing backing array |
| 2 | `DeepCopyMap(src)` | Copy map without sharing |
| 3 | `SafeDivideEx(a, b)` | Divide with recover from panic |
| 4 | `Flatten(matrix)` | 2D → 1D slice |
| 5 | `UniqueInts(nums)` | Remove duplicates (preserve order) |
| 6 | `ChunkSlice(s, n)` | Split slice into chunks of size n |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/stdlib/03_builtins/
```
