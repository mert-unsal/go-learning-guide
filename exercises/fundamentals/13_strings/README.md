# 📦 Module 13 — Strings: The Byte/Rune Duality

> **Topics covered:** 2-word string header · bytes vs runes vs characters · UTF-8 variable-width encoding · string immutability · backing array sharing · efficient concatenation
>
> **Deep dive:** [Chapter 03 — Strings: Immutability & Boxing](../../../learnings/03_strings_immutability_and_boxing.md)

---

## 🗺️ Learning Path

```
1. Read learnings/03_strings_immutability_and_boxing.md   ← How strings work under the hood
2. Run cmd/concepts/stdlib/01-strings-strconv/            ← See string operations in action
3. Open exercises.go                                      ← Implement the 12 exercises
4. Run go test -race -v ./...                             ← Make them all pass
```

---

## 📚 What You Will Learn

| Concept | Exercise | Under the Hood |
|---------|----------|---------------|
| `len()` returns bytes, not runes | Ex 1 | String header `{ptr, len}` — len field is byte count |
| Reversing multi-byte strings | Ex 2 | `[]rune` conversion decodes UTF-8 into fixed-width int32s |
| Rune indexing vs byte indexing | Ex 3 | `s[n]` gives a byte, not a character |
| ASCII detection | Ex 4 | UTF-8 property: ASCII chars are always single-byte (< 0x80) |
| Byte offsets in `for range` | Ex 5 | `for i, r := range s` — `i` jumps by rune byte-width |
| String immutability | Ex 6 | `[]byte(s)` copies backing bytes — string stays unchanged |
| Replacing runes in immutable strings | Ex 7 | Must rebuild: `[]rune` → modify → `string()` |
| Safe UTF-8 truncation | Ex 8 | `s[:n]` splits multi-byte runes — truncate by rune count instead |
| UTF-8 byte-width classes | Ex 9 | 1-4 bytes per rune: ASCII, accented, CJK, emoji |
| Efficient concatenation | Ex 10 | `strings.Builder` with `Grow()` — amortized O(1) vs O(n²) for `+=` |
| Substring memory leak | Ex 11 | Substrings share backing array — `strings.Clone` detaches |
| Case-insensitive comparison | Ex 12 | ASCII fold: `b \| 0x20` maps A-Z to a-z |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `ByteAndRuneCount(s)` | Return both byte length and rune count |
| 2 | `ReverseString(s)` | Reverse handling multi-byte runes correctly |
| 3 | `NthRune(s, n)` | Get the nth rune (not byte!) |
| 4 | `IsASCII(s)` | Check if all bytes are < 128 |
| 5 | `RuneByteOffsets(s)` | Byte offset where each rune starts |
| 6 | `ProveImmutability(s)` | **Demonstrate string ↔ []byte copy** |
| 7 | `ReplaceAtRuneIndex(s, idx, r)` | **Replace rune in immutable string** |
| 8 | `SafeTruncate(s, maxRunes)` | **Production pattern: truncate without breaking UTF-8** |
| 9 | `CountByteClasses(s)` | Count 1/2/3/4-byte UTF-8 sequences |
| 10 | `ConcatRepeat(s, n)` | **Must use strings.Builder** |
| 11 | `DetachSubstring(s, start, end)` | **Memory leak trap: strings.Clone** |
| 12 | `EqualFoldASCII(a, b)` | Case-insensitive compare (ASCII only) |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/fundamentals/13_strings/
```

---

## ✅ Done? Next Step

Move to **Phase 2 — Concurrency**:
```bash
go test -race -v ./exercises/fundamentals/10_goroutines/
```
