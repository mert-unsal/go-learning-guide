# 28 — reflect Under the Hood: TypeOf, ValueOf & When Reflection Is Justified

> **Companion exercises:** [exercises/stdlib/11_reflect](../exercises/stdlib/11_reflect/)

---

## Table of Contents

1. [The reflect Contract](#1-the-reflect-contract)
2. [reflect.Type — Compile-Time Metadata at Runtime](#2-reflecttype)
3. [reflect.Value — Runtime Values with Full Introspection](#3-reflectvalue)
4. [Struct Tags — The Metadata System](#4-struct-tags)
5. [Settability — Why You Need Pointers](#5-settability)
6. [The Nil Interface Trap](#6-the-nil-interface-trap)
7. [Performance Cost of Reflection](#7-performance-cost)
8. [When to Use Reflection (and When NOT To)](#8-when-to-use-reflection)
9. [Patterns in the Standard Library](#9-patterns-in-the-standard-library)
10. [Cost Table](#10-cost-table)
11. [Quick Reference Card](#11-quick-reference-card)
12. [Further Reading](#12-further-reading)

---

## 1. The reflect Contract

Go's `reflect` package exposes the runtime type system that the compiler
normally hides. It's built on three laws (from Rob Pike's blog post):

### Law 1: Reflection goes from interface value to reflection object

```go
var x float64 = 3.14
v := reflect.ValueOf(x)   // interface{} → reflect.Value
t := reflect.TypeOf(x)    // interface{} → reflect.Type
```

When you pass `x` to `reflect.ValueOf()`, Go boxes `x` into an `interface{}`
(since the parameter type is `any`). The reflect package then reads the
`eface` struct's type pointer and data pointer.

### Law 2: Reflection goes from reflection object to interface value

```go
y := v.Interface().(float64)  // reflect.Value → interface{} → float64
```

### Law 3: To modify a reflection object, the value must be settable

```go
v := reflect.ValueOf(x)    // v holds a COPY of x
v.SetFloat(7.1)            // PANIC: using unaddressable value

v = reflect.ValueOf(&x).Elem()  // v holds the original x
v.SetFloat(7.1)                  // OK: x is now 7.1
```

### Under the Hood: eface and iface

The `reflect` package reads the same `runtime.eface` and `runtime.iface`
structures we explored in chapter 06:

```
┌──────────────────────────────┐
│     eface (empty interface)  │
├──────────────────────────────┤
│  _type  *_type     ──────────│──► type descriptor (size, kind, methods)
│  data   unsafe.Pointer ─────│──► actual value data
└──────────────────────────────┘

reflect.TypeOf(v)  → reads _type pointer → wraps as reflect.Type
reflect.ValueOf(v) → reads both _type and data → wraps as reflect.Value
```

---

## 2. reflect.Type

`reflect.Type` is an interface that describes a Go type:

```go
t := reflect.TypeOf(Person{})
t.Name()      // "Person"
t.Kind()      // reflect.Struct
t.NumField()  // 4
t.Field(0)    // reflect.StructField{Name: "Name", Type: string, Tag: ...}
t.NumMethod() // number of exported methods
```

### Kind vs Name

| Value | Kind() | Name() |
|-------|--------|--------|
| `42` | `reflect.Int` | `"int"` |
| `type Score int; Score(42)` | `reflect.Int` | `"Score"` |
| `Person{}` | `reflect.Struct` | `"Person"` |
| `[]int{1,2}` | `reflect.Slice` | `""` (unnamed) |
| `map[string]int{}` | `reflect.Map` | `""` (unnamed) |

**Kind** is the underlying category (int, struct, slice, etc.).
**Name** is the defined type name (empty for unnamed composite types).

### Type Identity

```go
reflect.TypeOf(0) == reflect.TypeOf(0)           // true (same type)
reflect.TypeOf(int(0)) == reflect.TypeOf(int32(0))  // false (different types)
```

`reflect.Type` values can be compared with `==`. Two types are identical if
they have the same definition. Named types are distinct even if their
underlying types match.

### Getting Types Without Values

```go
// Get reflect.Type for an interface (use nil pointer trick)
errorType := reflect.TypeOf((*error)(nil)).Elem()
```

Since you can't create a value of an interface type directly, you use a typed
nil pointer and call `.Elem()` to get the interface type.

---

## 3. reflect.Value

`reflect.Value` wraps a runtime value with type metadata:

```go
v := reflect.ValueOf(42)
v.Kind()            // reflect.Int
v.Int()             // 42 (as int64)
v.Interface()       // any(42)
v.CanSet()          // false (it's a copy)
```

### Value Methods by Kind

| Kind | Read Methods | Write Methods |
|------|-------------|---------------|
| Int/Int64 | `v.Int()` | `v.SetInt(n)` |
| Float64 | `v.Float()` | `v.SetFloat(f)` |
| String | `v.String()` | `v.SetString(s)` |
| Bool | `v.Bool()` | `v.SetBool(b)` |
| Slice | `v.Len()`, `v.Index(i)` | `v.Set(other)` |
| Map | `v.MapKeys()`, `v.MapIndex(k)` | `v.SetMapIndex(k, v)` |
| Struct | `v.NumField()`, `v.Field(i)` | `v.Field(i).Set(v)` |
| Ptr | `v.Elem()` | `v.Set(other)` |
| Func | `v.Call(args)` | N/A |

### The Zero Value

```go
v := reflect.ValueOf(nil)  // v.IsValid() == false
v.Kind()                    // reflect.Invalid
```

Always check `v.IsValid()` before calling type-specific methods.

### Getting Addressable Values

```go
// WRONG: copy, not settable
v := reflect.ValueOf(myStruct)

// RIGHT: pointer's element is settable
v := reflect.ValueOf(&myStruct).Elem()
v.FieldByName("Name").SetString("new")  // ✅
```

---

## 4. Struct Tags

Struct tags are string literals attached to struct fields:

```go
type User struct {
    Name  string `json:"name" validate:"required" db:"user_name"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"-"`
}
```

### How Tags Work Under the Hood

Tags are stored in the compiled binary as part of the type descriptor.
The `reflect.StructField.Tag` field is a `reflect.StructTag` (just a string)
with a conventional format:

```
`key1:"value1" key2:"value2"`
```

The `Tag.Get("json")` method parses this string looking for the key. It's
string parsing at runtime, not a map lookup. But it's cached by the reflect
package after the first access.

### Standard Library Tag Keys

| Key | Used By | Example |
|-----|---------|---------|
| `json` | `encoding/json` | `json:"name,omitempty"` |
| `xml` | `encoding/xml` | `xml:"name,attr"` |
| `yaml` | `gopkg.in/yaml.v3` | `yaml:"name"` |
| `db` | `sqlx`, `gorm` | `db:"user_name"` |
| `validate` | `go-playground/validator` | `validate:"required,email"` |
| `mapstructure` | `mapstructure` | `mapstructure:"field_name"` |

### Tag Conventions

The `json` tag format is: `json:"name,option1,option2"`

| Option | Meaning |
|--------|---------|
| `omitempty` | Skip if zero value |
| `-` | Always skip this field |
| `string` | Encode number as JSON string |

---

## 5. Settability

This is where most reflect confusion happens.

### Why Settability Exists

When you call `reflect.ValueOf(x)`, the function receives a COPY of `x`
(Go is pass-by-value). Setting a field on a copy would be silently useless,
so Go panics instead.

```go
x := 42
v := reflect.ValueOf(x)
v.CanSet()  // false — v holds a copy of x

v = reflect.ValueOf(&x).Elem()
v.CanSet()  // true — v refers to the original x
v.SetInt(100)
// x is now 100
```

### Struct Field Settability

```go
type Person struct {
    Name    string  // exported → settable (if struct is addressable)
    private string  // unexported → NOT settable
}

p := Person{Name: "Go"}
v := reflect.ValueOf(&p).Elem()

v.FieldByName("Name").CanSet()    // true
v.FieldByName("private").CanSet() // false — unexported
```

### The Settability Chain

```
reflect.ValueOf(&s)           → pointer Value (not settable itself)
  .Elem()                     → struct Value (settable)
    .FieldByName("Name")      → field Value (settable if exported)
      .SetString("new")       → ✅ modifies original struct
```

---

## 6. The Nil Interface Trap

This is one of Go's most infamous gotchas:

```go
var p *MyError = nil
var err error = p

fmt.Println(err == nil)  // false!
```

Why? Because `err` is an `iface{tab: *itab_for_MyError, data: nil}`. The
interface has a non-nil type pointer, so the interface itself is not nil.

### Detecting True Nil with Reflect

```go
func IsNilSafe(v any) bool {
    if v == nil {
        return true
    }
    rv := reflect.ValueOf(v)
    switch rv.Kind() {
    case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
        return rv.IsNil()
    }
    return false
}
```

**Important:** Only certain kinds can be nil. Calling `IsNil()` on an int
Value panics.

---

## 7. Performance Cost

Reflection is significantly slower than direct code:

### Benchmarks (approximate)

| Operation | Direct | Reflect | Slowdown |
|-----------|--------|---------|----------|
| Read struct field | ~0.3ns | ~50ns | ~160x |
| Set struct field | ~0.3ns | ~100ns | ~330x |
| Call function | ~2ns | ~300ns | ~150x |
| Type assertion | ~1ns | ~30ns | ~30x |
| Create slice | ~5ns | ~200ns | ~40x |

### Why It's Slow

1. **No inlining**: reflect calls go through interface dispatch
2. **Runtime type checks**: every operation validates kind/type
3. **Allocation**: many reflect operations allocate (boxing values)
4. **No compiler optimizations**: the compiler can't optimize through reflect

### When the Cost Matters

- **Serialization hot paths** (encoding/json processes millions of values)
- **ORM query builders** (called per database query)
- **Validation frameworks** (called per HTTP request)

The `encoding/json` package mitigates this by caching reflected type info
in a `sync.Map` (see chapter 24). Custom marshalers bypass reflect entirely.

---

## 8. When to Use Reflection

### Justified Uses

1. **Serialization/deserialization**: encoding/json, encoding/xml, encoding/gob
2. **Struct tag processing**: ORMs, validators, config loaders
3. **Testing utilities**: `reflect.DeepEqual`, test assertion libraries
4. **Dependency injection frameworks**: (rare in Go, but exists)
5. **Generic-before-generics**: code that needed type flexibility before Go 1.18

### NOT Justified (Use Alternatives)

| Instead of Reflect | Use This |
|-------------------|----------|
| Generic containers | Generics (Go 1.18+) |
| Type switches on known types | `switch v := v.(type)` |
| Dynamic dispatch | Interfaces |
| Struct copying | Code generation |

### The Go Proverb

> *"Reflection is never clear."*

This means: reflection code is hard to read, hard to debug, and hard to
maintain. If there's a non-reflect solution, prefer it. Generics eliminated
80% of the cases where reflection was previously necessary.

---

## 9. Patterns in the Standard Library

### encoding/json

The entire `json.Marshal`/`json.Unmarshal` machinery is built on reflect:
- Walk struct fields with `reflect.Type.NumField()`
- Read struct tags for JSON key names
- Convert values to JSON via `reflect.Value.Interface()`
- Cache type info in `sync.Map` for performance

### fmt.Println

`fmt` uses reflect to determine how to format any value:
- Check for `Stringer` or `error` interfaces
- Switch on `Kind()` for formatting rules
- Recursively format struct fields

### database/sql

`Scan()` uses reflect to:
- Match column names to struct field names
- Convert database types to Go types dynamically

---

## 10. Cost Table

| Operation | Cost | Allocations | Notes |
|-----------|------|-------------|-------|
| `reflect.TypeOf(v)` | ~5ns | 0 | Reads eface, no alloc |
| `reflect.ValueOf(v)` | ~10ns | 0-1 | May box value types |
| `v.Kind()` | ~1ns | 0 | Field read |
| `v.FieldByName("X")` | ~50ns | 0 | Linear scan + string compare |
| `v.FieldByIndex([]int{0})` | ~5ns | 0 | Direct index |
| `v.SetString("x")` | ~100ns | 0 | Type check + set |
| `v.Call(args)` | ~300ns | 1+ | Allocates return values |
| `v.Interface()` | ~10ns | 0-1 | May allocate for boxing |
| `reflect.DeepEqual(a, b)` | ~500ns+ | varies | Recursive, allocation-heavy |

---

## 11. Quick Reference Card

```text
┌─────────────────────────────────────────────────────────────────┐
│                     REFLECT CHEAT SHEET                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Three Laws:                                                    │
│    1. interface → reflect object (TypeOf, ValueOf)              │
│    2. reflect object → interface (.Interface())                 │
│    3. Must be settable to modify (need pointer)                 │
│                                                                 │
│  Types:                                                         │
│    t := reflect.TypeOf(v)                                       │
│    t.Kind()          → underlying kind (Int, Struct, Ptr...)    │
│    t.Name()          → type name ("Person", "" for unnamed)     │
│    t.NumField()      → struct field count                       │
│    t.Field(i).Tag    → struct tag string                        │
│                                                                 │
│  Values:                                                        │
│    v := reflect.ValueOf(x)                                      │
│    v := reflect.ValueOf(&x).Elem()  ← settable!                │
│    v.FieldByName("X").SetString("y")                            │
│    v.Call(args)      → invoke function                          │
│    v.Interface()     → back to any                              │
│                                                                 │
│  Nil-safe check:                                                │
│    v == nil || (v.Kind() is nillable && v.IsNil())              │
│                                                                 │
│  Interface type without value:                                  │
│    reflect.TypeOf((*error)(nil)).Elem()                         │
│                                                                 │
│  Performance: 30-300x slower than direct code                   │
│  Use only when: serialization, tags, testing, dynamic types     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 12. Further Reading

- [The Laws of Reflection (Rob Pike)](https://go.dev/blog/laws-of-reflection) — the definitive article
- [`reflect` package source](https://cs.opensource.google/go/go/+/master:src/reflect/) — `value.go` and `type.go`
- [Go Data Structures: Interfaces (Russ Cox)](https://research.swtch.com/interfaces) — how iface/eface work
- [encoding/json source](https://cs.opensource.google/go/go/+/master:src/encoding/json/) — production reflect usage
- [go-playground/validator source](https://github.com/go-playground/validator) — tag-based validation

---

> **Companion exercises:** [exercises/stdlib/11_reflect](../exercises/stdlib/11_reflect/)
