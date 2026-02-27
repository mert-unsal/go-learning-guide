package packages_modules
import "testing"
func TestIsExported(t *testing.T) {
tests := []struct {
name string
want bool
}{
{"PublicFunc", true},
{"privateFunc", false},
{"MyType", true},
{"myType", false},
{"HTTP", true},
{"", false},
}
for _, tt := range tests {
got := IsExportedSolution(tt.name)
if got != tt.want {
t.Errorf("IsExported(%q) = %v, want %v", tt.name, got, tt.want)
}
}
}
func TestInitRan(t *testing.T) {
log := GetInitLog()
if len(log) == 0 {
t.Error("init() should have run and added to initLog")
}
if log[0] != "packages_modules.init ran" {
t.Errorf("initLog[0] = %q, want 'packages_modules.init ran'", log[0])
}
}
func TestBlankImportPurpose(t *testing.T) {
got := BlankImportPurposeSolution()
if got == "" {
t.Error("BlankImportPurpose should return a non-empty explanation")
}
}
func TestModulePath(t *testing.T) {
got := ModulePathSolution("github.com/user/myapp", "config")
want := "github.com/user/myapp/config"
if got != want {
t.Errorf("ModulePath = %q, want %q", got, want)
}
}