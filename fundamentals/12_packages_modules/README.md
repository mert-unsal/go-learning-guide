# 📦 Module 12 — Packages & Modules

> **Topics covered:** Exported vs unexported identifiers · `init()` functions · Blank import `_` · Build tags · Module paths

---

## 🗺️ Learning Path

```
1. Read concepts.go        ← Theory with runnable examples
2. Open exercises.go       ← Implement the TODO functions yourself
3. Run the tests below     ← Instant feedback on your code
4. Stuck? Open solutions.go ← Only after you have tried!
```

---

## 📚 What You Will Learn

| Concept | Where |
|---------|-------|
| Exported (`Uppercase`) vs unexported (`lowercase`) | Exercise 1 — `IsExported` |
| `init()` — auto-runs when package is imported | Exercise 2 — `GetInitLog` |
| Blank import `import _ "pkg"` | Exercise 3 — `BlankImportPurpose` |
| Build tags `//go:build` | Exercise 3 — `BuildTagPurpose` |
| Module path construction | Exercise 4 — `ModulePath` |
| `go.mod` structure | `concepts.go` |

---

## ✏️ Exercises

Open `exercises.go` and implement each function:

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `IsExported(name string) bool` | Return true if first char is A-Z |
| 2 | `GetInitLog() []string` | Already done — observe how `init()` works |
| 3a | `BlankImportPurpose() string` | Return the purpose of blank import |
| 3b | `BuildTagPurpose() string` | Return what `-tags` flag does |
| 4 | `ModulePath(moduleName, subPackage string) string` | Return `"module/subpackage"` |

---

## 🧪 Run Tests

> ⚠️ The `./fundamentals/...` paths work from the **project root** only.  
> If you are inside this folder, use `go test . -v` instead.

### From project root:
```bash
go test ./fundamentals/12_packages_modules/... -v
```

### From inside this folder:
```bash
go test . -v
```

### Run a single test (from inside this folder):
```bash
go test . -v -run TestIsExported
go test . -v -run TestInitLog
go test . -v -run TestBlankImport
go test . -v -run TestModulePath
```

---

## 💡 Key Hints

<details>
<summary>Exercise 1 — IsExported hint</summary>

In Go, a name is exported if its first character is an uppercase letter:
```go
func IsExported(name string) bool {
    if len(name) == 0 { return false }
    return name[0] >= 'A' && name[0] <= 'Z'
}
```
</details>

<details>
<summary>What is init() for?</summary>

`init()` runs **automatically** before `main()` (or before any code that imports the package). It's used for:
- Registering database drivers
- Setting up default configurations
- Validating environment variables at startup

A package can have multiple `init()` functions. They run in the order they appear in the file.
</details>

<details>
<summary>What is blank import `import _ "pkg"` for?</summary>

It imports a package **only for its side effects** (i.e., its `init()` functions run), without using any of its exports. Most common use case: registering a database driver.

```go
import _ "github.com/lib/pq"  // registers PostgreSQL driver via init()
```
</details>

---

## 🎉 You Finished All Fundamentals!

Run the full fundamentals test suite to confirm everything passes:
```bash
go test ./fundamentals/... -v
```

### Next steps — choose your path:
| Path | Command |
|------|---------|
| 🧩 LeetCode problems | `go test ./leetcode/... -v` |
| 🏆 HackerRank problems | `go test ./hackerrank/... -v` |
| 🔧 Practical Go (Docker, config, deploy) | Browse `practical/` folder |
| 📖 Stdlib deep-dive | Browse `stdlib/` folder |

