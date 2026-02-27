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
func IsExported(name string) bool {
// TODO: check if name[0] is uppercase (A-Z)
// Hint: name[0] >= 'A' && name[0] <= 'Z'
return false
}
// Exercise 2: init() order
// Go runs init() functions in the order of imports and file names.
// initLog records which init functions ran and in what order.
var initLog []string
func init() {
// This runs automatically when the package is imported.
// In real code: setup defaults, validate config, register drivers.
initLog = append(initLog, "packages_modules.init ran")
}
// GetInitLog returns the log of init calls (for testing).
func GetInitLog() []string {
return initLog
}
// Exercise 3: Blank identifier _ to suppress unused import warnings.
// Build tags let you include/exclude files by OS, arch, or custom tag.
// Answer the quiz by returning the correct string:
// Q: What does _ in "import _ fmt" do?
func BlankImportPurpose() string {
// TODO: return "runs init() functions of the package without using its exports"
return ""
}
// Q: What does go build -tags integration do?
func BuildTagPurpose() string {
// TODO: return "includes files with //go:build integration at the top"
return ""
}
// Exercise 4:
// Demonstrate module-aware path building.
// Given a module path and a sub-package name, return the full import path.
// Example: ModulePath("github.com/user/myapp", "config") → "github.com/user/myapp/config"
func ModulePath(moduleName, subPackage string) string {
// TODO: return moduleName + "/" + subPackage
return ""
}