# 📦 Module 04 — io & Files: Reader/Writer Composition

> **Topics covered:** io.Reader · io.Writer · LimitReader · MultiReader · TeeReader · io.Pipe · io.Copy · decorator pattern
>
> **Deep dive:** [Chapter 03 — Strings §5 (conversion costs)](../../../learnings/03_strings_immutability_and_boxing.md), [Chapter 19 §5 (decorator pattern)](../../../learnings/19_context_interface_masterclass.md)

---

## 🗺️ Learning Path

```
1. Read: io.Reader and io.Writer are 1-method interfaces        ← The foundation
2. Understand: Composition = wrapping readers/writers in layers  ← The power
3. Open exercises.go                                             ← Implement 12 exercises
4. Run go test -race -v ./...                                    ← Make them all pass
```

---

## 📚 What You Will Learn

| Concept | Exercise | Key Insight |
|---------|----------|-------------|
| File read/write round-trip | Ex 1 | `os.WriteFile` / `os.ReadFile` for simple cases |
| Line scanning | Ex 2 | `strings.NewReader` adapts string → `io.Reader` |
| Accept `io.Reader` (not concrete) | Ex 3 | Works with files, HTTP, strings, pipes — anything |
| `io.ReadAll` | Ex 4 | Reads until EOF — fine for small payloads |
| `io.LimitReader` | Ex 5 | **Cap untrusted input** (HTTP bodies, uploads) |
| `io.MultiReader` | Ex 6 | **Chain readers** end-to-end (decorator) |
| `io.TeeReader` | Ex 7 | **Read + capture** simultaneously (like Unix tee) |
| `io.Pipe` | Ex 8 | **Synchronous in-memory pipe** between goroutines |
| `io.Copy` | Ex 9 | **32KB streaming** — never loads all into memory |
| Reader→Filter→Writer | Ex 10 | **Pipeline pattern** for log/CSV/stream processing |
| Implement `io.Writer` | Ex 11 | **Decorator: add counting** without changing writer |
| `io.SectionReader` | Ex 12 | **Random access** — read from offset |

---

## ✏️ Exercises

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `WriteAndReadFile(name, content)` | File I/O round-trip |
| 2 | `CountLines(content)` | bufio.Scanner line counting |
| 3 | `WordFrequency(r)` | Word frequency from io.Reader |
| 4 | `CopyReader(src)` | io.ReadAll |
| 5 | `ReadLimited(r, max)` | **io.LimitReader** |
| 6 | `ConcatReaders(readers...)` | **io.MultiReader** |
| 7 | `ReadAndCapture(r)` | **io.TeeReader** |
| 8 | `PipeTransfer(msg)` | **io.Pipe + goroutine** |
| 9 | `CopyWithBuffer(src)` | **io.Copy streaming** |
| 10 | `FilterLines(r, w, keep)` | **Reader→Filter→Writer pipeline** |
| 11 | `WriterCounter.Write(p)` | **Implement io.Writer** |
| 12 | `SectionRead(r, offset, len)` | **io.SectionReader** |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/stdlib/04_io_files/
```

---

## ✅ Done? Next Step

```bash
go test -race -v ./exercises/stdlib/05_encoding_json/
```
