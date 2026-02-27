// Package dependency_management demonstrates Go module and dependency management.
//
// ============================================================
// GO MODULES — COMPLETE GUIDE
// ============================================================
//
// A Go module is a collection of packages with a go.mod file at the root.
// go.mod defines: module name, Go version, and required dependencies.
//
// ─────────────────────────────────────────────────────────────
// INITIALISE A NEW MODULE
// ─────────────────────────────────────────────────────────────
//   mkdir myapp && cd myapp
//   go mod init github.com/yourname/myapp
//
//   This creates go.mod:
//     module github.com/yourname/myapp
//     go 1.25.7
//
// ─────────────────────────────────────────────────────────────
// ADD A DEPENDENCY
// ─────────────────────────────────────────────────────────────
//   go get github.com/some/package           ← latest version
//   go get github.com/some/package@v1.2.3    ← specific version
//   go get github.com/some/package@latest    ← explicitly latest
//
//   After running go get:
//   • go.mod is updated with the require line
//   • go.sum is updated with cryptographic hashes (do NOT edit this)
//
// ─────────────────────────────────────────────────────────────
// COMMON DEPENDENCIES USED IN REAL PROJECTS
// ─────────────────────────────────────────────────────────────
//   go get github.com/joho/godotenv          ← .env file loading
//   go get gopkg.in/yaml.v3                  ← YAML parsing
//   go get github.com/gin-gonic/gin          ← HTTP web framework
//   go get go.uber.org/zap                   ← structured logging
//   go get github.com/stretchr/testify       ← test assertions
//   go get gorm.io/gorm                      ← ORM for databases
//
// ─────────────────────────────────────────────────────────────
// TIDY, VENDOR, VERIFY
// ─────────────────────────────────────────────────────────────
//   go mod tidy        ← remove unused deps, add missing ones
//   go mod vendor      ← copy all deps into ./vendor/ folder
//   go mod verify      ← check that deps haven't been tampered with
//   go mod download    ← pre-download all deps to local cache
//   go mod graph       ← print the dependency graph
//
// ─────────────────────────────────────────────────────────────
// UPGRADE / DOWNGRADE DEPENDENCIES
// ─────────────────────────────────────────────────────────────
//   go get github.com/some/package@v1.3.0    ← pin exact version
//   go get github.com/some/package@v1        ← latest v1.x.x
//   go get -u ./...                          ← upgrade ALL deps (careful!)
//   go get -u=patch ./...                    ← only patch version bumps
//
// ─────────────────────────────────────────────────────────────
// REMOVE A DEPENDENCY
// ─────────────────────────────────────────────────────────────
//   1. Remove all import references in your code
//   2. Run: go mod tidy
//
// ─────────────────────────────────────────────────────────────
// REPLACE DIRECTIVE (useful for local development)
// ─────────────────────────────────────────────────────────────
//   In go.mod:
//     replace github.com/some/package => ../local-fork
//
//   This lets you develop two modules side-by-side locally.
//
// ─────────────────────────────────────────────────────────────
// WORKSPACE MODE (Go 1.18+) — multiple modules together
// ─────────────────────────────────────────────────────────────
//   go work init ./moduleA ./moduleB
//   go work use ./moduleC
//
//   Creates go.work file — useful for monorepo development.
//   Do NOT commit go.work to source control (add to .gitignore).
//
// ─────────────────────────────────────────────────────────────
// PRIVATE MODULES
// ─────────────────────────────────────────────────────────────
//   export GONOSUMCHECK=github.com/yourprivate/*
//   export GOPRIVATE=github.com/yourprivate/*
//   export GONOSUMDB=github.com/yourprivate/*
//
// ─────────────────────────────────────────────────────────────
// go.sum FILE
// ─────────────────────────────────────────────────────────────
//   • Auto-generated — never edit manually
//   • Contains SHA-256 hashes of every dependency version
//   • MUST be committed to source control (ensures reproducibility)
//
// ─────────────────────────────────────────────────────────────
// READING go.mod EXAMPLE
// ─────────────────────────────────────────────────────────────
//   module github.com/yourname/myapp
//
//   go 1.25.7
//
//   require (
//       github.com/gin-gonic/gin v1.9.1
//       gopkg.in/yaml.v3 v3.0.1
//       github.com/joho/godotenv v1.5.1
//   )
//
//   require (
//       // indirect dependencies (transitive)
//       golang.org/x/net v0.17.0 // indirect
//   )

package dependency_management
