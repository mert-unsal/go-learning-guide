# 11 reflect — Runtime Type Introspection

> **Companion chapter:** [learnings/28_reflect_under_the_hood.md](../../../learnings/28_reflect_under_the_hood.md)

## Exercises

| # | Function | Concepts | Difficulty |
|---|---------|----------|------------|
| 1 | `TypeName` | `reflect.TypeOf`, Kind, Elem | ⭐ |
| 2 | `IsNilSafe` | Nil interface trap, `IsNil()` | ⭐⭐ |
| 3 | `FieldNames` | Struct field iteration, `NumField` | ⭐⭐ |
| 4 | `GetTag` | Struct tags, `Tag.Get()` | ⭐⭐ |
| 5 | `SetField` | `reflect.Value.Set`, settability | ⭐⭐⭐ |
| 6 | `StructToMap` | JSON tags, exported field check | ⭐⭐⭐ |
| 7 | `MapToStruct` | Populate struct from map dynamically | ⭐⭐⭐ |
| 8 | `CallFunc` | `reflect.Value.Call`, function invocation | ⭐⭐⭐ |
| 9 | `DeepEqualNilSafe` | `reflect.DeepEqual`, nil-safe comparison | ⭐⭐⭐ |
| 10 | `MakeSlice` | `reflect.MakeSlice`, `SliceOf` | ⭐⭐⭐ |
| 11 | `ImplementsError` | Interface check via `Implements` | ⭐⭐⭐ |
| 12 | `ValidateRequired` | Tag-driven validation, `IsZero()` | ⭐⭐⭐ |

## How to Practice

```bash
go test -race -v ./exercises/stdlib/11_reflect/
go test -race -run TestStructToMap ./exercises/stdlib/11_reflect/
```

## Key Insights

- **`reflect.TypeOf(v)`** returns static type info (fields, methods, tags)
- **`reflect.ValueOf(v)`** returns the runtime value (can read and set)
- **Settability** requires a pointer: `reflect.ValueOf(&s).Elem().FieldByName("X").Set(...)`
- **The nil interface trap**: `var err error = (*MyError)(nil)` is NOT nil
- **Struct tags** are compile-time metadata read at runtime via reflection
- *"Reflection is never clear"* — use it only when compile-time types are unknown
