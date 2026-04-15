# 📦 Module 05 — encoding/json: Marshal, Decode & Production Patterns

> **Topics covered:** struct tags · Decoder/Encoder streaming · custom MarshalJSON · json.RawMessage · json.Number · null vs absent
>
> **Deep dive:** [Chapter 24 — encoding/json Under the Hood](../../../learnings/24_encoding_json_under_the_hood.md)

---

## 🗺️ Learning Path

```
1. Read: Chapter 24 — encoding/json internals          ← How reflection drives marshal
2. Open exercises.go                                    ← Implement 12 exercises
3. Run go test -race -v ./...                           ← Make them all pass
```

---

## 📚 What You Will Learn

| Concept | Exercise | Key Insight |
|---------|----------|-------------|
| Marshal/Unmarshal round-trip | Ex 1 | `json.Marshal` → `json.Unmarshal` — the basics |
| Dynamic JSON with `map[string]any` | Ex 2 | Numbers become `float64` — the precision trap |
| Marshal slices | Ex 3 | Go slices → JSON arrays |
| Struct tags: `omitempty`, `"-"` | Ex 4 | **Control what gets serialized** |
| `json.NewDecoder` streaming | Ex 5 | **Read from io.Reader** — no ReadAll needed |
| `json.NewEncoder` streaming | Ex 6 | **Write to io.Writer** — direct output |
| Custom `MarshalJSON`/`UnmarshalJSON` | Ex 7 | **Override encoding for any type** |
| `json.RawMessage` delayed parsing | Ex 8 | **Parse envelope first, payload later** |
| Null vs absent fields | Ex 9 | **Pointer fields + raw check for true null** |
| Custom time format | Ex 10 | **Override time.Time's default RFC3339** |
| `json.Valid` | Ex 11 | **Validate without full decode** |
| `json.Number` precision | Ex 12 | **Preserve int64 beyond 2^53** |

---

## ✏️ Exercises

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `RoundTrip(e)` | Marshal → Unmarshal round-trip |
| 2 | `ParseJSON(s)` | Decode into `map[string]any` |
| 3 | `MarshalSlice(nums)` | `[]int` → JSON array string |
| 4 | `MarshalAppConfig(c)` | Struct tags in action |
| 5 | `DecodeStream(r)` | `json.NewDecoder` streaming |
| 6 | `EncodeStream(w, e)` | `json.NewEncoder` streaming |
| 7 | `StatusCode.MarshalJSON/UnmarshalJSON` | Custom marshal interface |
| 8 | `DelayParse(s)` | `json.RawMessage` pattern |
| 9 | `NullableFields(s)` | Null vs absent detection |
| 10 | `Event.MarshalJSON/UnmarshalJSON` | Custom time format |
| 11 | `ValidateJSON(s)` | `json.Valid` |
| 12 | `PreciseNumber(s)` | `json.Number` for large ints |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/stdlib/05_encoding_json/
```

---

## ✅ Done? Next Step

```bash
go test -race -v ./exercises/stdlib/06_math/
```
