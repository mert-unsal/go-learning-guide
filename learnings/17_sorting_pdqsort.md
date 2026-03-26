# pdqsort — Pattern-Defeating Quicksort

> The algorithm Go's `slices.Sort` (and `sort.Slice`) uses under the hood since Go 1.19.

---

## First, why not just use classic Quicksort?

Classic Quicksort is fast on average — **O(n log n)** — but has a nasty worst case:

```
Already sorted input:  [1, 2, 3, 4, 5, 6, 7, 8]
Classic Quicksort with bad pivot choice → O(n²) 💀
```

Real-world data is often **nearly sorted**, has **many duplicates**, or follows some pattern. Classic Quicksort falls apart on these.

---

## pdqsort = 3 algorithms stitched together smartly

pdqsort is not one algorithm. It **detects the shape of your data** and picks the best tool for it.

```
┌─────────────────────────────────────────────────────┐
│                     pdqsort                         │
│                                                     │
│   Is the slice tiny?  ──► Insertion Sort  O(n²)    │
│                            (but tiny n, so fast)    │
│                                                     │
│   Is data nearly sorted? ──► HeapSort  O(n log n)  │
│                               (bad pivot detected)  │
│                                                     │
│   Otherwise ──────────────► Quicksort  O(n log n)  │
│                               (with smart pivot)    │
└─────────────────────────────────────────────────────┘
```

---

## The 3 ingredients

### 1. Insertion Sort (for tiny slices, n ≤ 12)

```
[5, 3, 1, 4, 2]

Take each element, insert it into the right place:
[3, 5, 1, 4, 2]  → picked 3, inserted before 5
[1, 3, 5, 4, 2]  → picked 1, inserted at front
[1, 3, 4, 5, 2]  → picked 4, inserted after 3
[1, 2, 3, 4, 5]  → picked 2, inserted after 1
```

**Why use it for small slices?**
Insertion sort has zero overhead (no recursion, no pivot logic). For n ≤ 12, it beats everything.

---

### 2. Quicksort with Median-of-3 pivot (the main engine)

The key insight: **pivot choice determines everything**.

```
Bad pivot (always pick first):         Good pivot (median of 3):

[1, 2, 3, 4, 5]                        candidates: first=1, mid=3, last=5
pivot = 1                               pivot = 3  ✅ (the middle value)
left  = []                              left  = [1, 2]
right = [2, 3, 4, 5]                    right = [4, 5]
→ perfectly unbalanced 💀               → balanced split ✅
```

**Median of 3**: look at the first, middle, and last elements. Use the median as the pivot. This avoids the worst case on already-sorted data.

```go
// Concept:
first, mid, last := data[0], data[n/2], data[n-1]
pivot := median(first, mid, last)
```

---

### 3. HeapSort as the escape hatch (when Quicksort misbehaves)

pdqsort **counts bad pivot choices**. If it detects too many (threshold: log₂n), it switches to HeapSort.

```
Normal run:     Quicksort → Quicksort → Quicksort → done  ✅
Pathological:   Quicksort → bad pivot → bad pivot → SWITCH TO HEAPSORT 🔄
```

HeapSort is always **O(n log n)** — never worse. It's slower than Quicksort on average, but it's the safety net that **guarantees** pdqsort never degrades to O(n²).

---

## The "Pattern-Defeating" part

The name comes from its ability to **detect and exploit common patterns** in real data.

### Pattern: Already sorted (or reverse sorted)

```
[1, 2, 3, 4, 5]  ← ascending run detected
```

pdqsort checks if the data is already sorted (or reverse sorted) before doing anything. If yes → done or reverse in O(n).

### Pattern: Many duplicates

```
[3, 3, 3, 1, 3, 3, 2, 3]
```

Uses a **3-way partition** (Dutch National Flag):
```
Left: elements < pivot  | Middle: elements == pivot  | Right: elements > pivot
[1, 2]                  | [3, 3, 3, 3, 3]            | []
```

The middle "equal" section is never recursed into again — huge speedup for data with many duplicates.

---

## The full decision tree

```
pdqsort(slice)
│
├── len ≤ 12?
│   └── InsertionSort → done
│
├── Check for sorted / reverse-sorted runs
│   └── already sorted? → done (or reverse) → done
│
├── Pick pivot (median of 3)
│
├── 3-way partition:
│   ├── [< pivot] → recurse
│   ├── [= pivot] → skip (already in place)
│   └── [> pivot] → recurse
│
└── Too many bad pivots? (counted internally)
    └── Switch to HeapSort → guaranteed O(n log n)
```

---

## Complexity Summary

| Scenario | Time Complexity | Algorithm used |
|---|---|---|
| Tiny slice (n ≤ 12) | O(n²) but tiny | Insertion Sort |
| Average case | O(n log n) | Quicksort |
| Already sorted | O(n) | Early exit |
| Many duplicates | O(n log n) | 3-way Quicksort |
| Pathological input | O(n log n) **guaranteed** | HeapSort fallback |
| Space | O(log n) | Stack frames only |

---

## Why Go chose pdqsort

Before Go 1.19, Go used a hand-rolled introsort (intro = introspective). pdqsort replaced it because:

- **Faster on real data** — pattern detection avoids unnecessary work
- **Same worst-case guarantee** — O(n log n) always
- **Cache friendly** — insertion sort on small slices fits in CPU cache
- **Simpler to reason about** — clear separation of responsibilities

---

## Key Takeaways

```
1. pdqsort = Quicksort + InsertionSort + HeapSort, combined intelligently

2. It DETECTS the shape of your data and adapts:
   - Tiny? → InsertionSort
   - Sorted? → Early exit
   - Duplicates? → 3-way partition
   - Quicksort going bad? → HeapSort

3. Worst case is ALWAYS O(n log n) — never O(n²)

4. This is what slices.Sort() and sort.Slice() use in Go today
```

