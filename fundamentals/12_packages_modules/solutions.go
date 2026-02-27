package packages_modules
// ============================================================
// SOLUTIONS — 12 Packages & Modules
// ============================================================
func IsExportedSolution(name string) bool {
if len(name) == 0 {
return false
}
return name[0] >= 'A' && name[0] <= 'Z'
}
func BlankImportPurposeSolution() string {
return "runs init() functions of the package without using its exports"
}
func BuildTagPurposeSolution() string {
return "includes files with //go:build integration at the top"
}
func ModulePathSolution(moduleName, subPackage string) string {
return moduleName + "/" + subPackage
}