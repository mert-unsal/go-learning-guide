// Package Concepts — standalone demonstration of Go's package system,
// exports, init(), modules, imports, internal packages, and build tags.
//
// Run: go run ./cmd/concepts/packages-modules/01-package-concepts
package main

import "fmt"

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
// This module:        import "gointerviewprep/utils"
// Third party:        import "github.com/stretchr/testify/assert"
//
// Import aliases (avoid name conflicts or long paths):
//   import (
//       "fmt"
//       myfmt "gointerviewprep/utils"   // use as myfmt.Something()
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
	fmt.Println("=== Package Concepts ===")
	fmt.Println(PublicFunction())
	fmt.Println("PublicConstant:", PublicConstant)
	fmt.Println("privateConstant:", privateConstant) // accessible here (same package)
}
