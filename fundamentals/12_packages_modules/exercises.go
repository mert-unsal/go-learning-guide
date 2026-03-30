package packages_modules

// ============================================================
// EXERCISES — 12 Packages & Modules
// ============================================================

// These exercises are more conceptual — you answer by implementing
// functions that demonstrate understanding of visibility, init, and imports.

// Exercise 1: Visibility rules
// exported (public) = starts with Uppercase
// unexported (private) = starts with lowercase
// ExportedValue is visible from other packages.
const ExportedValue = "I am exported"

// unexportedValue is only visible within this package.
const unexportedValue = "I am unexported" //nolint

// IsExported returns true if the first character of name is uppercase.
//
// LESSON: In Go, there are no public/private keywords.
// Capitalization IS the access modifier. This is enforced by the compiler.
// Any identifier starting with a capital letter is exported (visible from other packages).
func IsExported(name string) bool {
	return false
}

// Exercise 2: init() order
// Go runs init() functions in the order of imports and file names.
// initLog records which init functions ran and in what order.
var initLog []string

func init() {
	// This runs automatically when the package is imported.
	// In real code: setup defaults, validate config, register drivers.
}

// GetInitLog returns the log of init calls (for testing).
func GetInitLog() []string {
	return nil
}

// Exercise 3: Blank identifier _ to suppress unused import warnings.
// Build tags let you include/exclude files by OS, arch, or custom tag.

// BlankImportPurpose explains what `import _ "pkg"` does.
//
// LESSON: A blank import runs the package's init() functions without
// making its exports available. Used to register drivers (e.g., database drivers,
// image format decoders) that self-register via init().
// Example: `import _ "github.com/lib/pq"` registers the PostgreSQL driver.
func BlankImportPurpose() string {
	return ""
}

// BuildTagPurpose explains what `go build -tags integration` does.
//
// LESSON: Build tags (//go:build integration) let you selectively compile files.
// Common uses: integration tests, OS-specific code, feature flags.
func BuildTagPurpose() string {
	return ""
}

// Exercise 4:
// ModulePath builds a full import path from a module root and sub-package.
//
// LESSON: Go module paths are just slash-joined strings.
// The go.mod file defines the module root. Every sub-package's import path
// is: moduleName + "/" + relative/path/to/package.
func ModulePath(moduleName, subPackage string) string {
	return ""
}
