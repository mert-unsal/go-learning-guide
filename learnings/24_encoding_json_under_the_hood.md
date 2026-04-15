# 24 — encoding/json Under the Hood

> **JSON is Go's most-used encoding — and the one most engineers use wrong.**
> Understanding how `encoding/json` works at the reflection/compiler level,
> what it costs, and where the production traps are, separates a Go developer
> who writes correct code from one who writes fast, correct code.

---

## Table of Contents

1. [The Architecture: Reflect-Driven Encoding](#1-the-architecture-reflect-driven-encoding)
2. [Marshal: From Struct to Bytes](#2-marshal-from-struct-to-bytes)
3. [Struct Tags: The Compiler Metadata System](#3-struct-tags-the-compiler-metadata-system)
4. [Unmarshal: From Bytes to Struct](#4-unmarshal-from-bytes-to-struct)
5. [Decoder vs Unmarshal: Streaming vs Buffered](#5-decoder-vs-unmarshal-streaming-vs-buffered)
6. [Custom Marshalers: The Interface Hooks](#6-custom-marshalers-the-interface-hooks)
7. [json.RawMessage: Delayed Parsing](#7-jsonrawmessage-delayed-parsing)
8. [The Number Precision Trap](#8-the-number-precision-trap)
9. [Null vs Absent: The Pointer Pattern](#9-null-vs-absent-the-pointer-pattern)
10. [Performance: What Costs What](#10-performance-what-costs-what)
11. [Production Patterns](#11-production-patterns)
12. [Performance Cost Table](#12-performance-cost-table)
13. [Quick Reference Card](#13-quick-reference-card)
14. [Further Reading](#14-further-reading)

---

## 1. The Architecture: Reflect-Driven Encoding

`encoding/json` is built entirely on Go's `reflect` package. When you call
`json.Marshal(v)`, the package inspects `v` at runtime using reflection to
determine its type, fields, and tags.

```
┌───────────────────────────────────────────────────────────┐
│  json.Marshal(v) call chain                               │
│                                                           │
│  1. reflect.TypeOf(v) → get the type                     │
│  2. Look up type in encoderCache (sync.Map)              │
│     HIT  → use cached encoder function                   │
│     MISS → build encoder via newTypeEncoder(t)           │
│  3. Encoder function writes JSON to internal buffer      │
│  4. Return buffer contents as []byte                     │
│                                                           │
│  The encoder function is type-specific:                   │
│  - structEncoder: iterates fields, calls field encoders  │
│  - sliceEncoder: iterates elements                       │
│  - mapEncoder: iterates key-value pairs                  │
│  - ptrEncoder: dereference, delegate to element encoder  │
│  - interfaceEncoder: extract concrete type, recurse      │
└───────────────────────────────────────────────────────────┘
```

### The Encoder Cache

The first time you marshal a type, `encoding/json` builds an encoder function
for that specific type and caches it in a package-level `sync.Map`. Subsequent
marshals of the same type skip the reflection setup and use the cached function.

```go
// encoding/json/encode.go — simplified
var encoderCache sync.Map // map[reflect.Type]encoderFunc

func typeEncoder(t reflect.Type) encoderFunc {
    if fi, ok := encoderCache.Load(t); ok {
        return fi.(encoderFunc)
    }
    // Build encoder, cache it
    f := newTypeEncoder(t, true)
    encoderCache.Store(t, f)
    return f
}
```

**Production implication:** The first marshal of a new type is slow (reflection
+ function construction). All subsequent marshals are fast (cached function).
In a long-running server, this cost is amortized to near zero.

**Comparison with other languages:**
- **Java (Jackson):** Also reflection-based, but can use compile-time annotation
  processing for code generation. Has `ObjectMapper` instance caching
- **C# (System.Text.Json):** Source generators (compile-time) available since .NET 6,
  avoiding reflection entirely
- **Go:** Reflection only (no source generators in stdlib). Third-party libraries
  like `easyjson`, `go-json`, and `sonic` use code generation

---

## 2. Marshal: From Struct to Bytes

### The Struct Encoder

When marshaling a struct, `encoding/json` reflects over its fields:

```go
type Employee struct {
    Name       string `json:"name"`
    Department string `json:"department"`
    Salary     int    `json:"salary,omitempty"`
}
```

The generated `structEncoder` for this type effectively does:

```
1. Write '{'
2. For each exported field (in declaration order):
   a. Read the struct tag → determine JSON key name and options
   b. If omitempty AND field is zero value → skip
   c. Write '"key":' 
   d. Call the field's type encoder (stringEncoder, intEncoder, etc.)
3. Write '}'
```

### Field Discovery (structFields)

`encoding/json` uses `reflect.Type.NumField()` and `reflect.Type.Field(i)` to
discover struct fields. For each field, it:

1. Checks if the field is exported (must be — unexported fields are invisible)
2. Parses the `json` struct tag
3. Handles embedded structs (promoted fields)
4. Resolves conflicts (two fields with the same JSON name at different depths)

This field list is computed once per type and cached.

### The Internal Buffer

Marshal writes to an internal `encodeState` struct that wraps a `bytes.Buffer`:

```go
// encoding/json/encode.go — simplified
type encodeState struct {
    bytes.Buffer
    scratch [64]byte // avoid allocation for small writes
}
```

The `scratch` array avoids heap allocation for formatting small numbers and
short strings. This is a common optimization in the Go stdlib.

---

## 3. Struct Tags: The Compiler Metadata System

Struct tags are string literals attached to struct fields. They're stored in
the binary's type metadata and accessible via `reflect.StructField.Tag`.

```go
type Config struct {
    Host   string `json:"host"`                    // rename to "host"
    Port   int    `json:"port"`                    // rename to "port"
    Debug  bool   `json:"debug,omitempty"`         // omit if false
    Secret string `json:"-"`                       // never marshal
    Label  string `json:"label,string"`            // wrap number/bool in quotes
    Extra  string `json:",omitempty"`              // keep Go field name, add omitempty
}
```

### Tag Syntax

```
`json:"fieldname,opt1,opt2"`
```

| Tag | Meaning |
|-----|---------|
| `json:"name"` | Use "name" as JSON key |
| `json:"-"` | Skip this field entirely |
| `json:",omitempty"` | Omit if zero value (keep Go name) |
| `json:"name,omitempty"` | Custom name + omit if zero |
| `json:",string"` | Marshal/unmarshal as JSON string (for numbers/bools) |
| `json:"-,"` | Use literal "-" as JSON key (rare escape hatch) |

### What "Zero Value" Means for omitempty

| Type | Zero value (omitted) |
|------|---------------------|
| `bool` | `false` |
| `int`, `float64`, etc. | `0` |
| `string` | `""` |
| `pointer` | `nil` |
| `slice` | `nil` (NOT empty slice `[]int{}`) |
| `map` | `nil` (NOT empty map `map[string]int{}`) |
| `struct` | **Never omitted** — even if all fields are zero |
| `array` | **Never omitted** — fixed-size, always present |
| `interface` | `nil` |

**Gotcha:** `omitempty` does NOT omit empty structs. A struct with all-zero
fields is still serialized. This is because the zero value of a struct is
well-defined but checking "are all fields zero" would require deep reflection.

**Gotcha:** An empty slice (`[]int{}`) is NOT omitted — only a nil slice is.
`make([]int, 0)` produces an empty (non-nil) slice that serializes as `[]`.

---

## 4. Unmarshal: From Bytes to Struct

### The Scanner

`json.Unmarshal` first validates and tokenizes the JSON using a finite state
machine scanner (`encoding/json/scanner.go`). The scanner processes byte by
byte, tracking nesting depth and the current state.

```
┌───────────────────────────────────────────────────────────┐
│  json.Unmarshal(data, &v)                                │
│                                                           │
│  1. Scanner validates JSON syntax (byte-by-byte FSM)     │
│  2. reflect.TypeOf(v) → must be pointer to settable type │
│  3. Look up field mapping (cached per type)              │
│  4. Walk JSON tokens + reflect fields simultaneously:    │
│     - JSON key → find matching struct field              │
│     - JSON value → decode into field's type              │
│  5. Return any errors (unknown fields, type mismatches)  │
└───────────────────────────────────────────────────────────┘
```

### Field Matching

For each JSON key, unmarshal searches for a matching struct field:

1. **Exact match** on the JSON tag name
2. **Case-insensitive match** on the field name (if no tag)

This case-insensitive fallback is why `json:"Name"` matches `"name"` in JSON.
It uses `strings.EqualFold` internally — O(n) per lookup but cached per type.

### Unknown Fields

By default, unknown JSON keys are silently ignored. To reject them:

```go
dec := json.NewDecoder(r)
dec.DisallowUnknownFields()
err := dec.Decode(&v)  // returns error if JSON has extra keys
```

**Production pattern:** Use `DisallowUnknownFields()` for API endpoints
where strict schema validation matters (internal services). Skip it for
external APIs where forward compatibility is needed.

---

## 5. Decoder vs Unmarshal: Streaming vs Buffered

### json.Unmarshal

```go
var v MyStruct
err := json.Unmarshal([]byte(data), &v)
```

- Requires the **entire JSON** as `[]byte` in memory
- Validates the entire input before returning
- Simple, safe, correct for small payloads

### json.Decoder

```go
dec := json.NewDecoder(reader)
err := dec.Decode(&v)
```

- Reads from `io.Reader` incrementally
- Does NOT read the entire input into memory
- Can decode multiple JSON values from the same stream
- Has configuration methods: `UseNumber()`, `DisallowUnknownFields()`

### When to Use Which

| Scenario | Use |
|----------|-----|
| Small payload already in memory | `json.Unmarshal` |
| HTTP request body | `json.NewDecoder(r.Body)` |
| File processing | `json.NewDecoder(file)` |
| Multiple JSON objects in stream (NDJSON) | `json.NewDecoder` in a loop |
| Need `UseNumber()` or `DisallowUnknownFields()` | `json.NewDecoder` (only way) |

### Decoder for NDJSON (Newline-Delimited JSON)

```go
dec := json.NewDecoder(reader)
for dec.More() {
    var v Event
    if err := dec.Decode(&v); err != nil {
        return err
    }
    process(v)
}
```

`dec.More()` returns true if there are more values in the stream.
This pattern processes log files, event streams, and Kafka messages.

### json.Encoder

```go
enc := json.NewEncoder(writer)
enc.SetIndent("", "  ")  // pretty-print
err := enc.Encode(v)      // writes JSON + newline
```

**Gotcha:** `Encode` appends a `\n` after the JSON. This is intentional
(for streaming), but surprising if you're building an HTTP response:

```go
// This adds a trailing newline to the HTTP body
json.NewEncoder(w).Encode(data)

// If you don't want the newline:
data, _ := json.Marshal(v)
w.Write(data)
```

---

## 6. Custom Marshalers: The Interface Hooks

`encoding/json` checks two interfaces before using reflection:

```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```

If your type implements these, `encoding/json` calls your method instead of
using reflection. This is how you customize serialization.

### Implementation Pattern

```go
type StatusCode int

const (
    StatusOK    StatusCode = 200
    StatusNotFound StatusCode = 404
)

var statusNames = map[StatusCode]string{
    200: "OK",
    404: "Not Found",
    500: "Internal Server Error",
}

var statusCodes = map[string]StatusCode{
    "OK":                    200,
    "Not Found":             404,
    "Internal Server Error": 500,
}

func (s StatusCode) MarshalJSON() ([]byte, error) {
    name, ok := statusNames[s]
    if !ok {
        return nil, fmt.Errorf("unknown status code: %d", s)
    }
    return json.Marshal(name) // returns `"OK"` (with quotes)
}

func (s *StatusCode) UnmarshalJSON(data []byte) error {
    var name string
    if err := json.Unmarshal(data, &name); err != nil {
        return err
    }
    code, ok := statusCodes[name]
    if !ok {
        return fmt.Errorf("unknown status name: %s", name)
    }
    *s = code
    return nil
}
```

### The Receiver Rule

- `MarshalJSON` can be on value or pointer receiver (both work)
- `UnmarshalJSON` MUST be on pointer receiver (`*T`) — it needs to modify `*s`
- If you put `UnmarshalJSON` on value receiver, `json.Unmarshal` won't find it
  when decoding into `*T`

This follows from Go's method set rules: a value of type `T` has methods
with value receivers. A value of type `*T` has methods with both value and
pointer receivers. `json.Unmarshal` always works with `*T`.

### time.Time's Custom Marshaler

`time.Time` implements `MarshalJSON` to produce RFC3339:

```go
// time/time.go
func (t Time) MarshalJSON() ([]byte, error) {
    b := make([]byte, 0, len(RFC3339Nano)+len(`""`))
    b = append(b, '"')
    b = t.AppendFormat(b, RFC3339Nano)
    b = append(b, '"')
    return b, nil
}
```

To use a custom format, wrap the type or implement `MarshalJSON` on your struct.

---

## 7. json.RawMessage: Delayed Parsing

`json.RawMessage` is just `[]byte` with `MarshalJSON`/`UnmarshalJSON` that
pass the raw bytes through unchanged:

```go
type RawMessage []byte

func (m RawMessage) MarshalJSON() ([]byte, error) {
    return m, nil
}

func (m *RawMessage) UnmarshalJSON(data []byte) error {
    *m = append((*m)[0:0], data...)
    return nil
}
```

### The Envelope Pattern

The primary use case is "parse the envelope, then decide how to parse the payload":

```go
type Message struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"` // stays as raw bytes
}

func decode(input []byte) error {
    var msg Message
    json.Unmarshal(input, &msg)
    
    switch msg.Type {
    case "employee":
        var emp Employee
        json.Unmarshal(msg.Data, &emp) // now parse the data
    case "config":
        var cfg Config
        json.Unmarshal(msg.Data, &cfg)
    }
}
```

Without `RawMessage`, you'd have to:
1. Parse into `map[string]any` (lose type safety)
2. Or parse twice (once for type, once for data)
3. Or use a union struct with all possible fields (brittle)

**Production use:** API gateways, event processors, message queues where the
outer envelope is standard but the payload varies by type.

---

## 8. The Number Precision Trap

By default, `json.Unmarshal` into `any`/`interface{}` turns JSON numbers
into `float64`. IEEE 754 float64 has 53 bits of mantissa, which means
integers larger than 2^53 (9,007,199,254,740,992) lose precision:

```go
var m map[string]any
json.Unmarshal([]byte(`{"id":9007199254740993}`), &m)
fmt.Println(m["id"]) // 9.007199254740992e+15 — WRONG! Off by 1
```

### The Fix: json.Number

```go
dec := json.NewDecoder(strings.NewReader(`{"id":9007199254740993}`))
dec.UseNumber() // preserve exact number strings
var m map[string]any
dec.Decode(&m)

num := m["id"].(json.Number) // type is json.Number, not float64
str := num.String()           // "9007199254740993" — exact
val, _ := num.Int64()         // 9007199254740993 — correct
```

`json.Number` is just `type Number string` — it holds the raw JSON string
representation and provides `Int64()`, `Float64()`, and `String()` methods.

**When this matters:**
- Database IDs (snowflake IDs, Twitter IDs)
- Blockchain values
- Financial amounts (though you should use strings for those)
- Any int64 value that might exceed 2^53

**Production pattern:** If your API receives arbitrary JSON, always use
`json.NewDecoder` with `UseNumber()` instead of `json.Unmarshal` into `any`.

---

## 9. Null vs Absent: The Pointer Pattern

One of the most subtle JSON handling problems: distinguishing between a
field that's explicitly `null` and one that's simply absent.

```json
{"name": "Alice", "nickname": null}     // nickname is explicitly null
{"name": "Alice"}                        // nickname is absent
{"name": "Alice", "nickname": "Ali"}    // nickname has a value
```

### Pointer Fields

Using `*string` instead of `string` lets you detect null vs value:

```go
type Profile struct {
    Name     string  `json:"name"`
    Nickname *string `json:"nickname,omitempty"`
}
```

After unmarshal:
- Value present: `Nickname` points to the string
- Null: `Nickname` is `nil`
- Absent: `Nickname` is `nil`

**Problem:** Null and absent both result in `nil`. You can't distinguish them
with pointer fields alone.

### Distinguishing Null from Absent

To truly distinguish, you need to check the raw JSON:

```go
func isExplicitNull(jsonStr string, field string) bool {
    var raw map[string]json.RawMessage
    json.Unmarshal([]byte(jsonStr), &raw)
    
    val, exists := raw[field]
    if !exists {
        return false // absent
    }
    return string(val) == "null" // explicit null
}
```

**Production pattern:** For PATCH APIs where "set to null" and "don't change"
are different operations, use `json.RawMessage` or a custom unmarshaler.

---

## 10. Performance: What Costs What

### Allocation Profile

```go
// json.Marshal allocates:
// 1. encodeState (internal buffer) — pooled via sync.Pool
// 2. Result []byte — the output
// 3. Intermediate strings for map keys (if map)
// Total: ~2-3 allocations for a simple struct

// json.Unmarshal allocates:
// 1. scanner state
// 2. Any string fields (copied from input)
// 3. Any slice/map fields (created fresh)
// Total: 1 + N allocations (N = number of string/slice/map fields)
```

### The sync.Pool for encodeState

`encoding/json` uses `sync.Pool` to reuse `encodeState` buffers:

```go
// encoding/json/encode.go
var encodeStatePool sync.Pool

func newEncodeState() *encodeState {
    if v := encodeStatePool.Get(); v != nil {
        e := v.(*encodeState)
        e.Reset()
        return e
    }
    return new(encodeState)
}
```

This means repeated `json.Marshal` calls reuse buffers, reducing GC pressure.

### Benchmark Comparison

Typical numbers for a 10-field struct (Go 1.22, amd64):

| Operation | Time | Allocs |
|-----------|------|--------|
| `json.Marshal` (small struct) | ~500ns | 2 |
| `json.Unmarshal` (small struct) | ~800ns | 3-5 |
| `json.Marshal` (100-element slice) | ~15μs | 2 |
| `json.Unmarshal` (100-element slice) | ~25μs | ~102 |
| `json.NewEncoder.Encode` (to Writer) | ~450ns | 1 |
| `json.NewDecoder.Decode` (from Reader) | ~900ns | 4-6 |

### Faster Alternatives

| Library | Approach | Speedup | Trade-off |
|---------|----------|---------|-----------|
| `encoding/json` (stdlib) | Reflection | 1x (baseline) | Always available, stable |
| `github.com/goccy/go-json` | Optimized reflection | ~2-3x | Drop-in replacement |
| `github.com/bytedance/sonic` | JIT + SIMD | ~5-10x | Linux amd64 only, CGO |
| `github.com/mailru/easyjson` | Code generation | ~3-5x | Requires `go generate` |
| `github.com/json-iterator/go` | Reflection + codegen | ~2-3x | API-compatible |

**Production advice:** Start with `encoding/json`. Profile. If JSON is truly
the bottleneck (it rarely is), try `go-json` (drop-in, no CGO). Only reach
for `sonic` or `easyjson` if profiling proves JSON is >20% of your CPU time.

---

## 11. Production Patterns

### Pattern 1: HTTP Handler with Proper Error Handling

```go
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Limit body size
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
    
    // 2. Decode with strict mode
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    
    var req CreateUserRequest
    if err := dec.Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    
    // 3. Reject multiple JSON values
    if dec.More() {
        http.Error(w, "request body must contain a single JSON object",
            http.StatusBadRequest)
        return
    }
    
    // 4. Validate domain rules
    if err := req.Validate(); err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }
    
    // 5. Process and respond
    user := processRequest(req)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

### Pattern 2: Envelope Responses

```go
type Response struct {
    Data  any    `json:"data,omitempty"`
    Error string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{Data: data})
}

func writeError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{Error: msg})
}
```

### Pattern 3: Config File Loading

```go
func LoadConfig(path string) (*Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open config: %w", err)
    }
    defer f.Close()
    
    var cfg Config
    dec := json.NewDecoder(f)
    dec.DisallowUnknownFields() // catch typos in config
    if err := dec.Decode(&cfg); err != nil {
        return nil, fmt.Errorf("decode config: %w", err)
    }
    return &cfg, nil
}
```

### Pattern 4: NDJSON Stream Processing

```go
func ProcessEvents(r io.Reader) error {
    dec := json.NewDecoder(r)
    for dec.More() {
        var event Event
        if err := dec.Decode(&event); err != nil {
            return fmt.Errorf("decode event: %w", err)
        }
        if err := handleEvent(event); err != nil {
            return fmt.Errorf("handle event: %w", err)
        }
    }
    return nil
}
```

---

## 12. Performance Cost Table

| Operation | Cost | Allocs | Notes |
|-----------|------|--------|-------|
| `json.Marshal` (small struct) | ~500ns | 2 | Cached encoder after first call |
| `json.Marshal` (large struct, 50 fields) | ~5μs | 2 | Reflection cost scales linearly |
| `json.Unmarshal` (small struct) | ~800ns | 3-5 | String fields each allocate |
| `json.NewDecoder.Decode` | ~900ns | 4-6 | Slightly more than Unmarshal |
| Encoder cache lookup | ~50ns | 0 | `sync.Map.Load` |
| First marshal of new type | ~10μs | 10+ | Builds + caches encoder |
| Custom `MarshalJSON` | varies | varies | Bypasses reflection entirely |
| `json.Valid` | ~200ns/KB | 0 | Scanner only, no decode |
| String field unmarshal | ~100ns | 1 | Copies bytes to new string |
| `json.RawMessage` unmarshal | ~50ns | 1 | Just copies raw bytes |

---

## 13. Quick Reference Card

```
┌───────────────────────────────────────────────────────────────────┐
│  encoding/json — QUICK REFERENCE                                 │
├───────────────────────────────────────────────────────────────────┤
│                                                                   │
│  MARSHAL / UNMARSHAL                                              │
│  json.Marshal(v)           → []byte, error                       │
│  json.MarshalIndent(v,p,i) → []byte, error  (pretty-print)      │
│  json.Unmarshal(data, &v)  → error                               │
│  json.Valid(data)          → bool            (syntax check only) │
│                                                                   │
│  STREAMING                                                        │
│  json.NewDecoder(r)        Decode from io.Reader                 │
│  json.NewEncoder(w)        Encode to io.Writer                   │
│  dec.UseNumber()           Preserve exact numbers                │
│  dec.DisallowUnknownFields() Reject extra keys                   │
│  dec.More()                More values in stream?                │
│  enc.SetIndent(p, i)       Pretty-print output                  │
│  enc.SetEscapeHTML(false)  Don't escape <, >, & in strings      │
│                                                                   │
│  STRUCT TAGS                                                      │
│  `json:"name"`             Rename field                          │
│  `json:"-"`                Exclude field                         │
│  `json:",omitempty"`       Omit zero value                       │
│  `json:",string"`          Encode number/bool as JSON string     │
│  `json:"-,"`               Literal "-" as key name               │
│                                                                   │
│  CUSTOM ENCODING                                                  │
│  json.Marshaler            MarshalJSON() ([]byte, error)         │
│  json.Unmarshaler          UnmarshalJSON([]byte) error           │
│  json.RawMessage           Delay parsing (raw []byte)            │
│  json.Number               Preserve number precision (string)    │
│                                                                   │
│  OMITEMPTY ZERO VALUES                                            │
│  bool=false  int=0  string=""  pointer=nil  slice=nil  map=nil   │
│  struct=NEVER  array=NEVER  []T{}=NOT nil (serializes as [])     │
│                                                                   │
│  PRODUCTION RULES                                                 │
│  1. Limit HTTP body: http.MaxBytesReader(w, r.Body, max)        │
│  2. Strict API: dec.DisallowUnknownFields()                      │
│  3. Single object: check dec.More() after Decode                 │
│  4. Large int: dec.UseNumber() to avoid float64 truncation       │
│  5. Content-Type: always set "application/json"                  │
│  6. Encoder adds \n — use Marshal if you don't want it           │
│  7. Profile before switching to faster JSON libraries            │
└───────────────────────────────────────────────────────────────────┘
```

---

## 14. Further Reading

- **encoding/json source:** `go/src/encoding/json/encode.go` — Start with
  `Marshal()` and trace through `newTypeEncoder()`
- **encoding/json/scanner.go:** The FSM scanner — surprisingly readable
- **Go Blog: "JSON and Go"** — Official introduction to encoding/json
- **encoding/json/v2 proposal:** `github.com/go-json-experiment/json` —
  The next-generation JSON package being developed by the Go team with
  better defaults (case-sensitive matching, no trailing newline in Encoder)
- **sonic:** `github.com/bytedance/sonic` — JIT-based JSON for when stdlib
  is too slow. Understand the tradeoffs before adopting

---

## Companion Exercises

Practice these concepts:
→ [`exercises/stdlib/05_encoding_json/`](../exercises/stdlib/05_encoding_json/) — 12 exercises
covering struct tags, streaming decode/encode, custom marshalers, RawMessage,
null detection, custom time format, json.Valid, and json.Number precision.
