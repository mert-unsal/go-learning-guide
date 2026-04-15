# 02 Escape Analysis — Stack vs Heap Allocation

> **Companion chapter:** [learnings/13_memory_gc_escape_sorting.md](../../../learnings/13_memory_gc_escape_sorting.md)

## Exercises

| # | Function | Concepts | Difficulty |
|---|---------|----------|------------|
| 1 | `StackOnly` | No escape, pure stack | ⭐ |
| 2 | `EscapeToHeap` | Pointer return forces escape | ⭐ |
| 3 | `NoEscapeSlice` | Small slice stays on stack | ⭐⭐ |
| 4 | `EscapeSlice` | Returned slice escapes | ⭐⭐ |
| 5 | `InterfaceEscape` | Boxing forces escape | ⭐⭐ |
| 6 | `ClosureCapture` | Closure captures cause escape | ⭐⭐⭐ |
| 7 | `PreallocateVsAppend` | Pre-allocation avoids growslice | ⭐⭐ |
| 8 | `JoinStrings` | strings.Builder vs fmt.Sprintf | ⭐⭐ |
| 9 | `Distance` | Struct by value (no escape) | ⭐⭐ |
| 10 | `FormatRecord` | sync.Pool buffer reuse | ⭐⭐⭐ |
| 11 | `RemoveAt` | Slice trick, no allocation | ⭐⭐ |
| 12 | `BytesEqualString` | Zero-copy comparison | ⭐⭐⭐ |

## How to Practice

```bash
go test -race -v ./exercises/advanced/02_escape_analysis/

# THE KEY TOOL: see what escapes to heap
go build -gcflags='-m' ./exercises/advanced/02_escape_analysis/ 2>&1

# More verbose escape analysis
go build -gcflags='-m -m' ./exercises/advanced/02_escape_analysis/ 2>&1
```
