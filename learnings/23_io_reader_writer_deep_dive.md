# 23 — io.Reader / io.Writer: The Composable I/O Engine

> **Go's I/O system is one interface, one method, infinite composition.**
> Understanding how `io.Reader` and `io.Writer` work — and how the runtime,
> standard library, and your own code compose them — is the key to writing
> efficient, production-grade streaming code in Go.

---

## Table of Contents

1. [The One-Method Contract](#1-the-one-method-contract)
2. [io.Reader Under the Hood](#2-ioreader-under-the-hood)
3. [io.Writer Under the Hood](#3-iowriter-under-the-hood)
4. [The Decorator Pattern: Composition Over Inheritance](#4-the-decorator-pattern-composition-over-inheritance)
5. [Key Readers in the Standard Library](#5-key-readers-in-the-standard-library)
6. [Key Writers in the Standard Library](#6-key-writers-in-the-standard-library)
7. [io.Copy: The 32KB Streaming Engine](#7-iocopy-the-32kb-streaming-engine)
8. [io.Pipe: Synchronous In-Memory Bridge](#8-iopipe-synchronous-in-memory-bridge)
9. [bufio: The Buffering Layer](#9-bufio-the-buffering-layer)
10. [Zero-Copy Patterns and Optimizations](#10-zero-copy-patterns-and-optimizations)
11. [Production Patterns](#11-production-patterns)
12. [Performance Cost Table](#12-performance-cost-table)
13. [Quick Reference Card](#13-quick-reference-card)
14. [Further Reading](#14-further-reading)

---

## 1. The One-Method Contract

Go's entire I/O system is built on two interfaces:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

That's it. One method each. Every file, every network connection, every HTTP
body, every compression stream, every encryption layer, every buffer — they all
implement one or both of these interfaces.

### Why This Matters

In Java, you have `InputStream`, `OutputStream`, `BufferedInputStream`,
`InputStreamReader`, `FileInputStream`, `ByteArrayInputStream` — a deep
class hierarchy. In Go, there's no hierarchy. There's **composition**.

```
Java:  BufferedReader(InputStreamReader(FileInputStream(file)))
       └── 3 classes, inheritance chain, constructor ceremony

Go:    bufio.NewReader(os.Open(file))
       └── 2 interfaces, same contract, zero inheritance
```

The power comes from the contract's simplicity:

1. **The caller provides the buffer** (`p []byte`) — the Reader doesn't allocate
2. **Partial reads are normal** — `n` can be less than `len(p)`
3. **EOF is not an error** — it's a signal (`io.EOF`)
4. **n > 0 AND err != nil is valid** — process the bytes, then handle the error

### The Reader Contract in Detail

```go
// n = number of bytes read into p[0:n]
// err = nil means more data may be available
// err = io.EOF means stream ended (n may still be > 0!)
// err = other means something went wrong

func (r *myReader) Read(p []byte) (n int, err error)
```

**Critical rule:** A Reader may return `n > 0, io.EOF` in the same call.
This happens when the last chunk of data fits in the buffer. Callers MUST
process the `n` bytes before checking `err`:

```go
// ✅ Correct: process bytes, then check error
n, err := r.Read(buf)
result = append(result, buf[:n]...)  // use n bytes first
if err == io.EOF {
    break
}

// ❌ Wrong: check error first, lose last bytes
n, err := r.Read(buf)
if err != nil {
    break  // dropped the last n bytes!
}
```

This is specified in `io/io.go:82-96` and is the single most common
io.Reader bug in Go code.

---

## 2. io.Reader Under the Hood

### What Happens When You Call Read()

There's no magic runtime support for `io.Reader`. It's a plain interface.
The compiler generates a standard `iface` dispatch:

```
┌─────────────────────────────────────────────────────┐
│  iface for io.Reader                                │
│  ┌──────────────┐                                   │
│  │ tab *itab    │──► itab { inter: io.Reader,       │
│  │              │         _type: *os.File,           │
│  │              │         fun[0]: (*File).Read }     │
│  ├──────────────┤                                   │
│  │ data unsafe  │──► *os.File { fd: 3, name: ... }  │
│  │   .Pointer   │                                   │
│  └──────────────┘                                   │
└─────────────────────────────────────────────────────┘
```

When you call `r.Read(buf)`, the runtime:
1. Loads `itab.fun[0]` — the concrete `Read` method pointer
2. Makes an indirect call through that pointer
3. The concrete implementation does the actual work

**Performance implication:** interface method calls are not inlineable.
In extremely hot paths (millions of calls per second), this matters.
The compiler cannot see through the interface to optimize the callee.

### Common Concrete Readers

| Type | What `Read()` does | Allocation |
|------|-------------------|------------|
| `*os.File` | `syscall.Read(fd, p)` — kernel call | Zero (reads into caller's buffer) |
| `*strings.Reader` | `copy(p, s[r.i:])` — memory copy | Zero (reads from existing string) |
| `*bytes.Reader` | `copy(p, s[r.i:])` — memory copy | Zero (reads from existing slice) |
| `*bytes.Buffer` | `copy(p, buf[off:])` — memory copy | Zero (reads from existing buffer) |
| `*bufio.Reader` | Fills 4KB internal buffer, copies from it | 4KB buffer (once, on creation) |
| `*net.TCPConn` | `syscall.Read(fd, p)` via netpoll | Zero (epoll/kqueue wakeup) |
| `http.Request.Body` | Usually `*net.TCPConn` + chunked reader | Depends on transfer encoding |

### The Caller-Provides-Buffer Pattern

This is a fundamental design decision. Compare with other languages:

```
Java InputStream.read():     allocates and returns byte[] internally
Python file.read():          allocates and returns bytes object
Go io.Reader.Read(p):        reads INTO caller-provided p
```

The Go approach means:
- **You control allocation.** Reuse the same `[]byte` across calls
- **You control buffer size.** 4KB for network, 32KB for file, 1MB for bulk
- **Zero-copy is possible.** The buffer can be a memory-mapped region
- **GC pressure is minimized.** No per-read allocation

---

## 3. io.Writer Under the Hood

### The Writer Contract

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

The contract is stricter than Reader:
- **Must return `n == len(p)` on success** — partial writes ARE errors
- **Must NOT modify `p`** — not even temporarily
- **Must NOT retain `p`** — the caller may reuse the buffer immediately

That second rule (`p` must not be retained) is subtle but critical.
This is why `bytes.Buffer.Write()` copies the data internally — it can't
keep a reference to the caller's slice.

### Writer Dispatch

Same `iface` pattern as Reader. No special runtime support.

```go
var w io.Writer = os.Stdout

// Compiler generates:
//   itab.fun[0](data, p)
// Which resolves to:
//   (*os.File).Write(stdout, p)
// Which calls:
//   syscall.Write(fd=1, p)
```

### The StringWriter Optimization

`io.StringWriter` (added Go 1.12) avoids the `[]byte(s)` conversion:

```go
type StringWriter interface {
    WriteString(s string) (n int, err error)
}
```

`io.WriteString(w, s)` checks if `w` implements `StringWriter`.
If yes, calls `WriteString` directly. If no, converts `s` to `[]byte`.

```go
// io/io.go — simplified
func WriteString(w Writer, s string) (n int, err error) {
    if sw, ok := w.(StringWriter); ok {
        return sw.WriteString(s) // no allocation!
    }
    return w.Write([]byte(s))    // allocates
}
```

**Why this matters:** `string` → `[]byte` conversion allocates and copies
(because strings are immutable, `[]byte` is mutable — they can't share
backing memory safely without escape analysis proving the `[]byte` doesn't
escape). `WriteString` avoids this entirely.

`bytes.Buffer`, `os.File`, `bufio.Writer`, `strings.Builder` — they all
implement `StringWriter`.

---

## 4. The Decorator Pattern: Composition Over Inheritance

This is where Go's I/O system becomes powerful. Every reader/writer in the
stdlib is a decorator: it wraps another reader/writer and adds behavior.

```
┌─────────────────────────────────────────────────────────┐
│  Decorator Chain (reading a gzip JSON file)             │
│                                                         │
│  io.ReadAll(                                            │
│    json.NewDecoder(                                     │
│      gzip.NewReader(               ◄─── decompress     │
│        bufio.NewReader(            ◄─── buffer 4KB     │
│          os.Open("data.json.gz")   ◄─── syscall.Read   │
│        )                                                │
│      )                                                  │
│    )                                                    │
│  )                                                      │
│                                                         │
│  Data flow: disk → kernel → 4KB buf → gunzip → JSON    │
│  Each layer only knows about io.Reader. Zero coupling.  │
└─────────────────────────────────────────────────────────┘
```

This is the **Decorator pattern** — the same one from the Gang of Four book,
but without classes, without inheritance, without a `Decorator` base class.
In Go, you just wrap interfaces.

### How to Write Your Own Decorator

```go
// CountingReader counts bytes read through it
type CountingReader struct {
    R     io.Reader
    Count int64
}

func (cr *CountingReader) Read(p []byte) (int, error) {
    n, err := cr.R.Read(p)
    cr.Count += int64(n)
    return n, err
}
```

That's it. No base class. No override. No virtual dispatch table.
Just: wrap the interface, delegate, add behavior.

### Contrast with Java

```java
// Java: needs abstract class hierarchy
class CountingInputStream extends FilterInputStream {
    private long count;
    
    public CountingInputStream(InputStream in) {
        super(in);  // abstract base class ceremony
    }
    
    @Override
    public int read(byte[] b, int off, int len) throws IOException {
        int n = super.read(b, off, len);
        if (n > 0) count += n;
        return n;
    }
    // Also must override read() and read(byte[]) for correctness
}
```

Go's version: 8 lines, no inheritance, no base class, works with ANY
`io.Reader` — files, networks, strings, pipes, other decorators.

---

## 5. Key Readers in the Standard Library

### io.LimitReader — Cap Untrusted Input

```go
func LimitReader(r Reader, n int64) Reader
```

Returns a Reader that stops after `n` bytes. Under the hood, it's trivial:

```go
// io/io.go — simplified
type LimitedReader struct {
    R Reader
    N int64   // bytes remaining
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
    if l.N <= 0 {
        return 0, EOF
    }
    if int64(len(p)) > l.N {
        p = p[0:l.N]  // shrink the buffer to remaining bytes
    }
    n, err = l.R.Read(p)
    l.N -= int64(n)
    return
}
```

**Production use:** ALWAYS limit HTTP request bodies:

```go
// Without limit: attacker sends 10GB body, OOM kill
body, _ := io.ReadAll(r.Body)

// With limit: max 1MB
body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))

// Even better: http.MaxBytesReader returns 413 status
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
```

### io.MultiReader — Chain Readers Sequentially

```go
func MultiReader(readers ...Reader) Reader
```

Reads from each reader in sequence. When one returns EOF, moves to the next.

```go
// Prepend a header to a file stream — zero-copy
header := strings.NewReader("HEADER\n")
file, _ := os.Open("data.csv")
combined := io.MultiReader(header, file)
// combined reads "HEADER\n" then the file contents
```

Internally: `multiReader` holds a `[]Reader` and an index. Each `Read()` tries
the current reader. On EOF, advances the index. When all exhausted, returns EOF.

### io.TeeReader — Read and Copy Simultaneously

```go
func TeeReader(r Reader, w Writer) Reader
```

Everything read from the returned Reader is also written to `w`. Like `tee(1)`.

```go
// io/io.go — the entire implementation
type teeReader struct {
    r Reader
    w Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
    n, err = t.r.Read(p)
    if n > 0 {
        if n, err := t.w.Write(p[:n]); err != nil {
            return n, err
        }
    }
    return
}
```

**Production use:** capture a request body for logging while passing it through:

```go
var bodyLog bytes.Buffer
tee := io.TeeReader(req.Body, &bodyLog)
// decode from tee (consumes body AND fills bodyLog)
json.NewDecoder(tee).Decode(&payload)
// bodyLog now has the raw JSON for audit logging
```

### io.SectionReader — Random Access Windows

```go
func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader
```

Creates a view into a subset of a `ReaderAt`. The `SectionReader` implements
`Read`, `Seek`, and `ReadAt` — all confined to the `[off, off+n)` window.

**Why it exists:** `*os.File` implements `ReaderAt` via the `pread(2)` syscall,
which reads at an offset without changing the file position. This means multiple
goroutines can read different sections of the same file concurrently — no mutex,
no seeking, no interference.

```go
f, _ := os.Open("huge.dat")
// 4 goroutines read 4 quarters of the file concurrently
for i := 0; i < 4; i++ {
    sr := io.NewSectionReader(f, int64(i)*chunkSize, chunkSize)
    go processSection(sr)
}
```

This is how parallel file processing works in Go. The kernel handles the
concurrency via `pread(2)`.

---

## 6. Key Writers in the Standard Library

### bytes.Buffer — The Swiss Army Knife

```go
var buf bytes.Buffer
buf.WriteString("hello")
fmt.Fprintf(&buf, " world %d", 42)
result := buf.String()
```

Under the hood (`bytes/buffer.go`):
- Internal `[]byte` starts at 64 bytes, grows with `runtime.growslice` strategy
- `WriteByte`, `WriteRune`, `WriteString`, `Write` — all avoid allocation per call
- `ReadFrom` optimizes: reads directly into the internal buffer
- `Bytes()` returns a slice of the internal buffer (no copy)
- `String()` does copy (because `string` is immutable, `[]byte` is not)

**Gotcha:** `Reset()` keeps the allocated buffer for reuse. `Truncate(0)` does
the same. Neither frees memory. If you're done with a large buffer, let the GC
collect it — don't keep it in a `sync.Pool` if the peak size was anomalous.

### strings.Builder — Optimized String Construction

```go
var b strings.Builder
b.WriteString("hello")
b.WriteString(" world")
s := b.String()  // no copy! (since Go 1.10)
```

`strings.Builder` wraps a `[]byte` but its `String()` method uses
`unsafe.String` to convert without copying:

```go
// strings/builder.go
func (b *Builder) String() string {
    return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}
```

This is safe because `Builder` enforces that after `String()` is called,
the buffer is never modified (it panics if you try to write after calling
`String()` — actually it doesn't panic, but the pointer is shared, so the
`noescape` trick plus the `copyCheck` prevents corruption).

**When to use which:**
| Need | Use |
|------|-----|
| Build a string from parts | `strings.Builder` |
| Read/write bytes, reuse buffer | `bytes.Buffer` |
| Format into a writer | `fmt.Fprintf(w, ...)` |
| Hot path, known size | `make([]byte, 0, size)` + manual append |

### io.Discard — The /dev/null Writer

```go
var Discard Writer = devNull(0)
```

Writes succeed, bytes are thrown away. Useful for benchmarking and testing:

```go
// How fast can we encode without I/O?
io.Copy(io.Discard, myReader)
```

---

## 7. io.Copy: The 32KB Streaming Engine

`io.Copy` is the workhorse of Go I/O. It streams data from a Reader to a
Writer in chunks, without ever holding the entire payload in memory.

```go
func Copy(dst Writer, src Reader) (written int64, err error)
```

### The Implementation (io/io.go)

Simplified from the actual source:

```go
func copyBuffer(dst Writer, src Reader, buf []byte) (int64, error) {
    // Optimization 1: WriterTo — source can write directly
    if wt, ok := src.(WriterTo); ok {
        return wt.WriteTo(dst)
    }
    // Optimization 2: ReaderFrom — destination can read directly
    if rt, ok := dst.(ReaderFrom); ok {
        return rt.ReadFrom(src)
    }
    // Fallback: manual pump with 32KB buffer
    if buf == nil {
        size := 32 * 1024
        if l, ok := src.(*LimitedReader); ok && int64(size) > l.N {
            size = int(l.N) // don't allocate more than needed
        }
        buf = make([]byte, size)
    }
    for {
        nr, er := src.Read(buf)
        if nr > 0 {
            nw, ew := dst.Write(buf[0:nr])
            if nw < 0 || nr < nw {
                nw = 0
                if ew == nil {
                    ew = errInvalidWrite
                }
            }
            written += int64(nw)
            if ew != nil {
                return written, ew
            }
            if nr != nw {
                return written, ErrShortWrite
            }
        }
        if er != nil {
            if er != EOF {
                err = er
            }
            break
        }
    }
    return written, err
}
```

### The Three Code Paths

```
┌──────────────────────────────────────────────────────────┐
│  io.Copy decision tree                                   │
│                                                          │
│  1. src implements WriterTo?                             │
│     YES → src.WriteTo(dst)         ◄─ fastest path      │
│                                      (e.g., *os.File     │
│                                       uses sendfile(2))  │
│                                                          │
│  2. dst implements ReaderFrom?                           │
│     YES → dst.ReadFrom(src)        ◄─ second fastest     │
│                                      (e.g., *bytes.Buffer│
│                                       reads directly     │
│                                       into internal buf) │
│                                                          │
│  3. Neither?                                             │
│     Allocate 32KB buf              ◄─ generic fallback   │
│     Loop: Read(buf) → Write(buf)                         │
└──────────────────────────────────────────────────────────┘
```

### The sendfile(2) Optimization

When copying from `*os.File` to `*net.TCPConn`, Go can use the kernel's
`sendfile(2)` syscall — data goes from disk to network socket without
ever entering user space:

```
Without sendfile:   disk → kernel buf → user buf → kernel buf → network
With sendfile:      disk → kernel buf ──────────────────────── → network
                                       (zero-copy in kernel)
```

This happens transparently through the `WriterTo` / `ReaderFrom` interface
check in `io.Copy`. The `*os.File.WriteTo` method (or `*net.TCPConn.ReadFrom`)
detects the optimization opportunity and uses `sendfile(2)`.

**Production implication:** when serving static files over HTTP, ensure you
use `http.ServeFile` or `http.ServeContent` — they use `io.Copy` internally,
which triggers `sendfile(2)`. Don't read the file into memory first.

---

## 8. io.Pipe: Synchronous In-Memory Bridge

```go
func Pipe() (*PipeReader, *PipeWriter)
```

Creates a synchronous, in-memory pipe. Writes block until reads consume the
data, and vice versa. There is no internal buffering.

### Under the Hood (io/pipe.go)

The pipe uses a single `sync.Mutex` and two `sync.Cond` variables:

```
┌──────────────────────────────────────────────────────────┐
│  io.pipe internals                                       │
│                                                          │
│  ┌─────────────┐    ┌─────────────┐                     │
│  │ PipeReader   │    │ PipeWriter   │                    │
│  │  p *pipe ────│────│── p *pipe    │                    │
│  └─────────────┘    └─────────────┘                     │
│          │                  │                            │
│          └───────┬──────────┘                            │
│                  ▼                                       │
│          ┌─────────────┐                                 │
│          │   pipe       │                                │
│          │   wrMu  sync.Mutex (serialize writes)        │
│          │   wrCh  chan []byte (data channel)            │
│          │   rdCh  chan int    (read-done signal)        │
│          │   once  sync.Once  (close once)              │
│          │   done  chan struct{} (close signal)          │
│          │   rerr  onceError                            │
│          │   werr  onceError                            │
│          └─────────────┘                                 │
└──────────────────────────────────────────────────────────┘
```

Actually, since Go 1.21, `io.Pipe` uses channels internally (`wrCh`, `rdCh`)
instead of mutexes + conds. The writer sends the `[]byte` slice on `wrCh`,
the reader reads from it directly (zero-copy within the process), then signals
completion on `rdCh` with the byte count.

### When to Use Pipe

**Use case 1: Connect a WriterTo with a ReaderFrom**

```go
pr, pw := io.Pipe()
go func() {
    json.NewEncoder(pw).Encode(data) // writes JSON to pipe
    pw.Close()
}()
http.Post(url, "application/json", pr) // reads from pipe
```

Without the pipe, you'd need to `json.Marshal` into a `[]byte` first (allocating
the entire payload), then wrap it in `bytes.NewReader`. The pipe streams it.

**Use case 2: Process output of exec.Command**

```go
pr, pw := io.Pipe()
cmd := exec.Command("grep", "-r", "TODO", ".")
cmd.Stdout = pw
go func() {
    cmd.Run()
    pw.Close()
}()
scanner := bufio.NewScanner(pr)
for scanner.Scan() {
    process(scanner.Text())
}
```

### Pipe Gotchas

1. **Deadlock risk:** Write blocks until Read consumes. If writer and reader are
   in the same goroutine, instant deadlock. ALWAYS put one side in a goroutine.
2. **Close the writer:** The reader gets `io.EOF` only when the writer calls
   `Close()` (or `CloseWithError()`). Forgetting to close = goroutine leak.
3. **Error propagation:** `CloseWithError(err)` makes the other side see that
   error. Use it for cancellation: `pw.CloseWithError(ctx.Err())`.

---

## 9. bufio: The Buffering Layer

### bufio.Reader

```go
func NewReader(rd io.Reader) *Reader         // 4096 byte buffer
func NewReaderSize(rd io.Reader, size int) *Reader
```

Under the hood (`bufio/bufio.go`):

```
┌──────────────────────────────────────────────────────────┐
│  bufio.Reader                                            │
│                                                          │
│  ┌────────────────────────────────────────────────┐      │
│  │ buf  [4096]byte                                │      │
│  │      ┌──────┬─────────────────┬────────────┐   │      │
│  │      │ used │  buffered data  │   empty    │   │      │
│  │      └──────┴─────────────────┴────────────┘   │      │
│  │      0      r                 w           cap  │      │
│  │                                                │      │
│  │ r   int     // read position                   │      │
│  │ w   int     // write position                  │      │
│  │ rd  Reader  // underlying reader               │      │
│  │ err error   // sticky error                    │      │
│  └────────────────────────────────────────────────┘      │
│                                                          │
│  Read(p):                                                │
│    1. If buf has data (w > r): copy to p, advance r      │
│    2. If buf empty AND len(p) >= len(buf):                │
│       read directly into p (skip buffer entirely!)       │
│    3. Otherwise: fill buf from rd, then copy to p        │
└──────────────────────────────────────────────────────────┘
```

**Key insight:** When the caller's buffer is larger than the internal buffer,
`bufio.Reader` bypasses its own buffer and reads directly into the caller's
slice. This prevents double-copying for large reads.

### bufio.Scanner

`Scanner` is the idiomatic way to read lines (or tokens) from a Reader:

```go
scanner := bufio.NewScanner(reader)
for scanner.Scan() {
    line := scanner.Text()  // no trailing \n
}
if err := scanner.Err(); err != nil {
    // handle error (NOT io.EOF — that's normal)
}
```

Default buffer: 64KB. Default split function: `ScanLines`. The scanner
allocates a token buffer and grows it as needed (up to `MaxScanTokenSize`
= 1MB by default).

**Gotcha:** `scanner.Text()` returns a string that references the scanner's
internal buffer. The string is valid only until the next `Scan()` call.
If you need to keep it, copy it: `s := strings.Clone(scanner.Text())`.
Actually, `Text()` returns a newly allocated string each time (it calls
`string(s.token)` which copies), so it's safe. But `Bytes()` returns the
internal slice — that one IS invalidated on next `Scan()`.

### bufio.Writer

```go
func NewWriter(w io.Writer) *Writer         // 4096 byte buffer
func NewWriterSize(w io.Writer, size int) *Writer
```

Batches small writes into 4KB chunks before flushing to the underlying writer.
**Critical:** you MUST call `Flush()` when done, or data is lost:

```go
bw := bufio.NewWriter(file)
defer bw.Flush()  // MUST flush — or lose buffered data
bw.WriteString("hello")
bw.WriteString(" world")
// Flush sends "hello world" in one syscall instead of two
```

**Production pattern:** `bufio.Writer` reduces syscall count. Each `Write` to
`*os.File` is a `syscall.Write` — each one costs ~100ns on Linux. If you're
writing 1000 small strings, that's 1000 syscalls without buffering vs ~1 with.

---

## 10. Zero-Copy Patterns and Optimizations

### ReaderFrom / WriterTo Interface Checks

The most important optimization in Go's I/O system is the `ReadFrom`/`WriteTo`
check in `io.Copy`:

```go
// If the source can write to the destination directly
type WriterTo interface {
    WriteTo(w Writer) (n int64, err error)
}

// If the destination can read from the source directly
type ReaderFrom interface {
    ReadFrom(r Reader) (n int64, err error)
}
```

When a type implements one of these, `io.Copy` skips the 32KB buffer
entirely and lets the type handle the transfer with its own internal
buffer or kernel assistance.

**Who implements WriterTo:**
- `*bytes.Buffer` → writes its internal slice directly
- `*bytes.Reader` → writes its internal slice directly
- `*strings.Reader` → writes the string's backing array directly
- `*os.File` → uses `sendfile(2)` when writing to a socket

**Who implements ReaderFrom:**
- `*bytes.Buffer` → grows internal slice and reads directly into it
- `*bufio.Writer` → reads into internal buffer
- `*os.File` → uses `splice(2)` or `sendfile(2)` when reading from a pipe/socket
- `*net.TCPConn` → uses `splice(2)` for kernel-level transfer

### unsafe.String and unsafe.Slice (Go 1.20+)

For truly zero-copy string ↔ []byte conversion in controlled scenarios:

```go
// []byte → string without copying (DANGEROUS if slice is later modified)
s := unsafe.String(unsafe.SliceData(b), len(b))

// string → []byte without copying (DANGEROUS: modifying the result is UB)
b := unsafe.Slice(unsafe.StringData(s), len(s))
```

The `strings.Builder.String()` method uses this internally. You should almost
never need it in application code — but understanding that it exists helps you
know why `Builder.String()` is O(1) and not O(n).

### splice(2) — Kernel-Level Pipe Transfer

On Linux, when copying between two file descriptors (pipes, sockets, files),
Go can use `splice(2)` to move data entirely in kernel space:

```
Without splice:  fd1 → kernel → user buf → kernel → fd2
With splice:     fd1 → kernel pipe buffer → fd2  (never enters user space)
```

This is used by `io.Copy` when both sides are `*os.File` or `*net.TCPConn`.
The Go runtime's `net` package sets this up transparently.

---

## 11. Production Patterns

### Pattern 1: Streaming JSON Encode/Decode

```go
// ❌ Bad: marshal entire payload into memory
data, _ := json.Marshal(bigStruct)
resp, _ := http.Post(url, "application/json", bytes.NewReader(data))

// ✅ Good: stream through a pipe
pr, pw := io.Pipe()
go func() {
    json.NewEncoder(pw).Encode(bigStruct)
    pw.Close()
}()
resp, _ := http.Post(url, "application/json", pr)
```

For decoding, `json.NewDecoder(resp.Body).Decode(&v)` already streams —
it reads tokens from the reader incrementally.

### Pattern 2: Limit All Untrusted Input

```go
// In HTTP handlers:
const maxBody = 10 << 20 // 10MB
r.Body = http.MaxBytesReader(w, r.Body, maxBody)

// When reading from any untrusted source:
limited := io.LimitReader(untrustedReader, maxSize)
data, err := io.ReadAll(limited)
```

### Pattern 3: Reader/Writer as API Boundary

Design functions to accept `io.Reader`/`io.Writer` instead of concrete types:

```go
// ❌ Rigid: only works with files
func ProcessFile(path string) error

// ✅ Flexible: works with files, HTTP, strings, pipes, tests
func Process(r io.Reader) error
```

This is the single most impactful API design pattern in Go. Your function
works with `os.Stdin`, `http.Request.Body`, `strings.NewReader("test")`,
`bytes.Buffer`, `gzip.NewReader(...)`, `*exec.Cmd.Stdout` — anything.

### Pattern 4: io.ReadAll vs io.Copy

```go
// io.ReadAll: loads ENTIRE content into memory
// Use for: small payloads (<10MB), config files, test fixtures
data, err := io.ReadAll(resp.Body)

// io.Copy: streams 32KB at a time, constant memory
// Use for: large payloads, file downloads, proxying
n, err := io.Copy(outFile, resp.Body)
```

**Rule of thumb:** if you need the whole `[]byte`, use `ReadAll`. If you're
just moving data from A to B, use `Copy`.

### Pattern 5: Multi-Writer for Logging + Processing

```go
// Write to both file and stdout simultaneously
f, _ := os.Create("output.log")
mw := io.MultiWriter(f, os.Stdout)
fmt.Fprintf(mw, "Server started on port %d\n", port)
```

### Pattern 6: Close the Body

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()  // MUST close — or leak TCP connection

// Even if you don't read the body:
io.Copy(io.Discard, resp.Body)  // drain before close (reuse connection)
resp.Body.Close()
```

HTTP/1.1 connections are reused by the transport pool. If you don't read and
close the body, the connection can't be returned to the pool.

---

## 12. Performance Cost Table

| Operation | Cost | Allocations | Notes |
|-----------|------|-------------|-------|
| `io.ReadAll(small)` | O(n) | 1-3 (growing slice) | Fine for <10MB |
| `io.ReadAll(large)` | O(n) | Many (repeated grow) | Avoid for >10MB |
| `io.Copy` | O(n) | 1 (32KB buffer) | Constant memory |
| `io.Copy` with `WriterTo` | O(n) | 0 | Source manages transfer |
| `io.Copy` file→socket | O(n) | 0 | `sendfile(2)` — zero-copy |
| `bufio.NewReader` | O(1) | 1 (4KB buffer) | Amortizes syscalls |
| `bufio.Scanner` | O(n) | 1 per line (Text()) | Token buf grows to 1MB max |
| `strings.NewReader` | O(1) | 1 (Reader struct) | No copy of string |
| `bytes.NewReader` | O(1) | 1 (Reader struct) | No copy of slice |
| `io.LimitReader` | O(1) | 1 (LimitedReader) | Wrapper only |
| `io.MultiReader` | O(1) | 1 (slice of readers) | No data copy |
| `io.TeeReader` | O(1) | 1 (teeReader struct) | Writes on every read |
| `io.Pipe` | O(1) | 1 (pipe struct) | Synchronous, no buffer |
| `io.Discard.Write` | O(1) | 0 | Bytes thrown away |

### Syscall Costs (Linux, approximate)

| Operation | Cost |
|-----------|------|
| `read(fd, buf, n)` | ~100ns for cached file |
| `write(fd, buf, n)` | ~100ns for buffered file |
| `sendfile(fd, fd, off, n)` | ~50ns + DMA transfer |
| `splice(fd, fd, n)` | ~30ns setup + kernel copy |

---

## 13. Quick Reference Card

```
┌───────────────────────────────────────────────────────────────────┐
│  io PACKAGE — QUICK REFERENCE                                    │
├───────────────────────────────────────────────────────────────────┤
│                                                                   │
│  CORE INTERFACES                                                  │
│  io.Reader         Read(p []byte) (n int, err error)             │
│  io.Writer         Write(p []byte) (n int, err error)            │
│  io.Closer         Close() error                                 │
│  io.Seeker         Seek(offset int64, whence int) (int64, error) │
│  io.ReaderAt       ReadAt(p []byte, off int64) (n int, err error)│
│  io.WriterTo       WriteTo(w Writer) (n int64, err error)        │
│  io.ReaderFrom     ReadFrom(r Reader) (n int64, err error)       │
│  io.StringWriter   WriteString(s string) (n int, err error)      │
│                                                                   │
│  COMPOSED INTERFACES                                              │
│  io.ReadCloser     Reader + Closer                               │
│  io.WriteCloser    Writer + Closer                               │
│  io.ReadWriter     Reader + Writer                               │
│  io.ReadWriteCloser Reader + Writer + Closer                     │
│  io.ReadSeeker     Reader + Seeker                               │
│  io.ReadWriteSeeker Reader + Writer + Seeker                     │
│                                                                   │
│  READERS (decorators)                                             │
│  io.LimitReader(r, n)         Cap at n bytes                     │
│  io.MultiReader(r1, r2, ...)  Chain sequentially                 │
│  io.TeeReader(r, w)           Read + copy to w                   │
│  io.NewSectionReader(r, off, n) Window into ReaderAt             │
│  io.NopCloser(r)              Add no-op Close()                  │
│                                                                   │
│  WRITERS                                                          │
│  io.MultiWriter(w1, w2, ...)  Write to all simultaneously        │
│  io.Discard                   /dev/null                           │
│                                                                   │
│  TRANSFER                                                         │
│  io.Copy(dst, src)            Stream 32KB chunks                 │
│  io.CopyN(dst, src, n)       Copy exactly n bytes               │
│  io.CopyBuffer(dst, src, buf) Use provided buffer               │
│  io.ReadAll(r)                Read entire stream                 │
│  io.ReadFull(r, buf)          Read exactly len(buf) bytes        │
│  io.WriteString(w, s)         String→Writer (no alloc if impl)  │
│                                                                   │
│  PIPE                                                             │
│  io.Pipe()                    Synchronous in-memory pipe         │
│  pr.Read(p) / pw.Write(p)    Block until paired call             │
│  pw.Close() / pw.CloseWithError(err)                             │
│                                                                   │
│  SENTINELS                                                        │
│  io.EOF                       Normal end of stream               │
│  io.ErrUnexpectedEOF          Stream ended prematurely           │
│  io.ErrShortWrite             Write returned n < len(p)          │
│  io.ErrClosedPipe             Write to closed pipe               │
│                                                                   │
│  BUFIO                                                            │
│  bufio.NewReader(r)           4KB buffered reader                │
│  bufio.NewScanner(r)          Line/token scanner (64KB buf)      │
│  bufio.NewWriter(w)           4KB buffered writer                │
│  bw.Flush()                   MUST call when done!               │
│                                                                   │
│  ESCAPE ANALYSIS                                                  │
│  go build -gcflags='-m' — check if your Reader/Writer escapes   │
│  Interface values always escape to heap (can't inline dispatch)  │
│  In hot paths: consider concrete types to enable inlining        │
│                                                                   │
│  GOLDEN RULES                                                     │
│  1. Process n bytes BEFORE checking err (n>0, EOF possible)      │
│  2. Accept io.Reader, not concrete types — enables composition   │
│  3. LimitReader on ALL untrusted input                           │
│  4. io.Copy for streaming — never ReadAll large data             │
│  5. bufio.Writer MUST Flush() — or lose data                    │
│  6. Close HTTP response bodies — or leak connections             │
│  7. Pipe writer in goroutine — or deadlock                       │
└───────────────────────────────────────────────────────────────────┘
```

---

## 14. Further Reading

- **io package source:** `go/src/io/io.go` — Read the entire file. It's 650 lines
  and every line is worth understanding
- **bufio package source:** `go/src/bufio/bufio.go` — See how buffering wraps readers
- **os.File.ReadFrom:** `go/src/os/readfrom_linux.go` — The sendfile optimization
- **net.TCPConn.ReadFrom:** `go/src/net/tcpsock_posix.go` — splice optimization
- **io.Pipe rewrite (Go 1.21):** Changed from mutex+cond to channels internally
- **Rob Pike, "Go at Google":** The io.Reader/Writer design was inspired by Plan 9
- **Russ Cox, "io/fs design":** How io interfaces influenced the fs package design

---

## Companion Exercises

Practice these concepts:
→ [`exercises/stdlib/04_io_files/`](../exercises/stdlib/04_io_files/) — 12 exercises
covering file I/O, LimitReader, MultiReader, TeeReader, Pipe, Copy, FilterLines,
WriterCounter, and SectionReader.
