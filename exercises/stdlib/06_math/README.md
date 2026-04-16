# 📦 Module 06 — math: Numeric Operations

> **Topics covered:** math.Sqrt · math.Pow · bit manipulation · GCD/LCM · rounding · clamping

---

## ✏️ Exercises

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `Hypotenuse(a, b)` | √(a² + b²) using math.Sqrt |
| 2 | `IsPowerOfTwo(n)` | Bit manipulation: n & (n-1) == 0 |
| 3 | `ClampEx(val, min, max)` | Clamp value to range |
| 4 | `RoundToN(f, n)` | Round to n decimal places |
| 5 | `GCD(a, b)` | Greatest common divisor (Euclidean) |
| 6 | `LCM(a, b)` | Least common multiple via GCD |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/stdlib/06_math/
```

---

## 📖 Companion Chapter

For the deep-dive theory behind these exercises, read:

- [13 — Memory, GC, Escape Analysis & Sorting](../../../learnings/13_memory_gc_escape_sorting.md) — numeric precision, `math` package internals, sorting algorithms
