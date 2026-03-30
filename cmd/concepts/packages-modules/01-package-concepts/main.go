// Package Concepts — standalone demonstration of Go's package system,
// exports, init(), modules, imports, internal packages, and build tags.
//
// Run: go run ./cmd/concepts/packages-modules/01-package-concepts
package main

import (
	"fmt"
	"go/build"
	"os"
	"runtime"
	"strings"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

// ============================================================
// 1. PACKAGES — The Basic Unit of Code Organization
// ============================================================
// Every Go source file belongs to exactly one package.
// Package name = directory name (by convention).
// The 'main' package is special: it produces an executable.
//
// Rules:
//   - All files in a directory share the same package name
//   - You cannot have two packages in the same directory (except _test packages)
//   - Package names should be short, lowercase, single words
//
// Good:  package http, package json, package sort
// Bad:   package HTTPClient, package my_package

// ============================================================
// 2. EXPORTED vs UNEXPORTED IDENTIFIERS
// ============================================================
// Go uses CAPITALIZATION as the visibility rule — no public/private keywords.
//
//   Exported   (public):  starts with Uppercase letter → accessible from other packages
//   Unexported (private): starts with lowercase letter → only within this package
//
// In this standalone main package, both are accessible. But when writing
// library packages, only Exported identifiers form the public API.

// PublicConstant is exported — visible from any package that imports this one.
const PublicConstant = "I am exported"

// privateConstant is unexported — only visible within this package.
const privateConstant = "I am unexported"

// PublicFunction is exported.
func PublicFunction() string {
	return privateHelper() // can call unexported functions within same package
}

// privateHelper is unexported.
func privateHelper() string {
	return "called from inside the package"
}

// ============================================================
// 3. THE init() FUNCTION
// ============================================================
// Each package can define one (or more) init() functions.
// They run automatically when the package is initialized, BEFORE main().
//
// Order:
//   1. Package-level variables are initialized first (in declaration order)
//   2. init() functions run next
//   3. If package A imports package B, B's init() runs before A's
//
// Use init() sparingly — it makes code harder to test and reason about.
// Common uses: registering drivers (database/sql), setting defaults.

var globalVar = computeGlobal() // runs before init()

func computeGlobal() string {
	return "global computed"
}

func init() {
	// This runs automatically when the package is loaded, BEFORE main().
	// Uncomment the line below to see it in action:
	// fmt.Println("init() called — globalVar:", globalVar)
	_ = globalVar // use globalVar to show ordering
}

// ============================================================
// 4. go.mod — THE MODULE SYSTEM
// ============================================================
// A module is a collection of packages with a shared go.mod file.
// go.mod lives at the root of the module and declares:
//   - The module path (used as the import prefix)
//   - The Go version
//   - Direct dependencies (require directives)
//
// Example go.mod:
//
//   module github.com/you/myproject
//
//   go 1.21
//
//   require (
//       github.com/some/lib v1.2.3
//   )
//
// Key commands:
//   go mod init <module-path>  — create a new module
//   go get <package>           — add a dependency
//   go mod tidy                — remove unused dependencies, add missing ones
//   go mod download            — download all dependencies to local cache
//   go list -m all             — list all dependencies

// ============================================================
// 5. IMPORTING PACKAGES
// ============================================================
//
// Standard library:   import "fmt"
// This module:        import "go-learning-guide/utils"
// Third party:        import "github.com/stretchr/testify/assert"
//
// Import aliases (avoid name conflicts or long paths):
//   import (
//       "fmt"
//       myfmt "go-learning-guide/utils"   // use as myfmt.Something()
//       _ "database/driver"              // blank import: only runs init(), no direct use
//   )
//
// Dot import (import all exported names into current scope — avoid in production):
//   import . "math"
//   x := Sqrt(4)  // no need to write math.Sqrt

// ============================================================
// 6. INTERNAL PACKAGES
// ============================================================
// A package path containing the element "internal" can only be imported
// by code in the parent tree of the "internal" directory.
//
// Example:
//   myproject/
//     internal/
//       auth/        ← can only be imported by packages under myproject/
//     cmd/
//       server/      ← CAN import myproject/internal/auth
//     external/      ← CAN'T import myproject/internal/auth (compile error)
//
// The compiler enforces this — it's not just a convention, it's a
// hard boundary. Use internal/ to hide implementation details while
// sharing code across your own packages.

// ============================================================
// 7. BUILD TAGS
// ============================================================
// Build tags control which files are included in a build.
// Place them at the top of a file, BEFORE the package clause.
//

// ============================================================
// 8. SUMMARY
// ============================================================

func main() {
	fmt.Printf("%s%s══════════════════════════════════════════%s\n", bold, blue, reset)
	fmt.Printf("%s%s  Go Package & Module Concepts            %s\n", bold, blue, reset)
	fmt.Printf("%s%s══════════════════════════════════════════%s\n\n", bold, blue, reset)

	// --- Exported vs Unexported ---
	fmt.Printf("%s▸ Exported vs Unexported identifiers%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Go uses CAPITALIZATION as the visibility rule — no public/private keywords%s\n", green, reset)
	fmt.Printf("  PublicConstant  = %s%q%s  ← Uppercase first letter → exported\n", magenta, PublicConstant, reset)
	fmt.Printf("  privateConstant = %s%q%s  ← lowercase first letter → unexported\n", magenta, privateConstant, reset)
	fmt.Printf("  PublicFunction()= %s%q%s  ← calls privateHelper() internally\n", magenta, PublicFunction(), reset)
	fmt.Printf("  %s✔ Both accessible here because we're in the same package (main)%s\n", green, reset)
	fmt.Printf("  %s⚠ In library packages, only Exported identifiers form the public API%s\n", yellow, reset)

	// --- init() ordering ---
	fmt.Printf("\n%s▸ init() function & initialization order%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Order: package-level vars → init() → main()%s\n", green, reset)
	fmt.Printf("  globalVar = %s%q%s  ← computed BEFORE init() ran\n", magenta, globalVar, reset)
	fmt.Printf("  %s⚠ Use init() sparingly — it makes code harder to test and reason about%s\n", yellow, reset)
	fmt.Printf("  %s✔ Common uses: registering DB drivers (database/sql), setting defaults%s\n", green, reset)

	// --- Module system ---
	fmt.Printf("\n%s▸ Module system (go.mod)%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ A module is a collection of packages with a shared go.mod%s\n", green, reset)

	// Read go.mod to show module info dynamically
	if data, err := os.ReadFile("go.mod"); err == nil {
		lines := strings.SplitN(string(data), "\n", 3)
		for _, line := range lines[:2] {
			line = strings.TrimSpace(line)
			if line != "" {
				fmt.Printf("    %s%s%s\n", magenta, line, reset)
			}
		}
	}

	fmt.Printf("  %s✔ Key commands: go mod init, go get, go mod tidy, go mod download%s\n", green, reset)

	// --- Imports ---
	fmt.Printf("\n%s▸ Import types%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Standard lib:%s    import \"fmt\"\n", green, reset)
	fmt.Printf("  %s✔ This module:%s     import \"go-learning-guide/utils\"\n", green, reset)
	fmt.Printf("  %s✔ Third party:%s     import \"github.com/stretchr/testify\"\n", green, reset)
	fmt.Printf("  %s✔ Alias:%s           import myfmt \"go-learning-guide/utils\"\n", green, reset)
	fmt.Printf("  %s✔ Blank import:%s    import _ \"database/driver\"  (only runs init())\n", green, reset)
	fmt.Printf("  %s⚠ Dot import:%s      import . \"math\"  (avoid in production!)%s\n", yellow, reset, reset)

	// --- internal/ packages ---
	fmt.Printf("\n%s▸ internal/ packages%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Packages under 'internal/' can only be imported by sibling/parent packages%s\n", green, reset)
	fmt.Printf("  %s✔ The COMPILER enforces this — not just a convention, a hard boundary%s\n", green, reset)
	fmt.Printf("  %s✔ Use internal/ to hide implementation details from external consumers%s\n", green, reset)

	// --- Build tags ---
	fmt.Printf("\n%s▸ Build tags%s\n", cyan+bold, reset)
	fmt.Printf("  %s✔ Control which files are included in a build%s\n", green, reset)
	fmt.Printf("  %s✔ New syntax (Go 1.17+): //go:build linux && amd64%s\n", green, reset)
	fmt.Printf("  %s⚠ Old syntax (deprecated): // +build linux,amd64%s\n", yellow, reset)

	// --- Runtime info ---
	fmt.Printf("\n%s▸ Current environment%s\n", cyan+bold, reset)
	fmt.Printf("  Go version:  %s%s%s\n", magenta, runtime.Version(), reset)
	fmt.Printf("  GOOS/GOARCH: %s%s/%s%s\n", magenta, runtime.GOOS, runtime.GOARCH, reset)
	fmt.Printf("  GOPATH:      %s%s%s\n", magenta, build.Default.GOPATH, reset)
	fmt.Printf("  NumCPU:      %s%d%s  ← GOMAXPROCS defaults to this\n", magenta, runtime.NumCPU(), reset)
}
