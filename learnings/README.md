# 📚 Learnings

This folder contains personal learning notes and deep-dive explanations collected during the Go interview prep journey.

Each file focuses on a specific concept that came up as a question or caused confusion — written in plain language with examples and diagrams.

---

## Index

| # | File | Topic |
|---|------|-------|
| 01 | [closures_scopes_values.md](./01_closures_scopes_values.md) | Scopes, variable shadowing, closures capturing variables vs values, pass-by-value |
| 02 | [closure_loop_gotcha.md](./02_closure_loop_gotcha.md) | The classic closure-in-a-loop bug and how to fix it with `i := i` |
| 03 | [pdqsort.md](./03_pdqsort.md) | How Go's `slices.Sort` works — pdqsort: Quicksort + InsertionSort + HeapSort |
| 04 | [auto_deref_auto_address.md](./04_auto_deref_auto_address.md) | When Go silently inserts `*` and `&` — auto-dereference and auto-address explained |
| 05 | [interface_nil_trap.md](./05_interface_nil_trap.md) | Interface internals: the two-word pair, the nil pointer trap, and the three guards |

---

> 💡 Add a new file here whenever you learn something worth remembering.

