# 01 Generics — Type Parameters & Constraints

> **Companion chapter:** [learnings/26_generics_under_the_hood.md](../../../learnings/26_generics_under_the_hood.md)

## Exercises

| # | Function / Type | Concepts | Difficulty |
|---|----------------|----------|------------|
| 1 | `Min` | `cmp.Ordered`, basic type parameter | ⭐ |
| 2 | `Contains` | `comparable` constraint, `==` operator | ⭐ |
| 3 | `Map` | Two type parameters `[T, U]` | ⭐⭐ |
| 4 | `Filter` | Function as parameter, predicate | ⭐⭐ |
| 5 | `Reduce` | Accumulator pattern, fold | ⭐⭐ |
| 6 | `Keys` / `Values` | Map type parameters `[K comparable, V any]` | ⭐⭐ |
| 7 | `Stack[T]` | Generic struct, method receivers | ⭐⭐⭐ |
| 8 | `Pair[K,V].Swap` | Multi-type struct, type swapping | ⭐⭐ |
| 9 | `MaxBy` | Comparator function, panic on empty | ⭐⭐⭐ |
| 10 | `GroupBy` | Key function, map of slices | ⭐⭐⭐ |
| 11 | `Result[T]` | Ok/Err pattern, zero value | ⭐⭐⭐ |
| 12 | `Set[T]` | `map[T]struct{}`, Union, Intersection | ⭐⭐⭐ |

## How to Practice

```bash
# Run all tests (they all fail initially)
go test -race -v ./exercises/advanced/01_generics/

# Run a specific test
go test -race -run TestStack ./exercises/advanced/01_generics/

# Check escape analysis for your generic code
go build -gcflags='-m' ./exercises/advanced/01_generics/ 2>&1 | head -30
```

## Key Insights

- **`any` vs `comparable` vs `cmp.Ordered`** — each constraint unlocks different operators
- **Generic structs** need type params on the receiver: `func (s *Stack[T]) Push(v T)`
- **`map[T]struct{}`** is the idiomatic Go set (struct{} costs zero bytes)
- **Two type params** `[T, U]` enable input/output type transformations
- Under the hood, Go uses **GC shape stenciling**: pointer types share one compiled version,
  value types each get their own. Use `go build -gcflags='-m'` to see instantiation decisions
