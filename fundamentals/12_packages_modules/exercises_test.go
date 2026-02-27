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
		got := IsExported(tt.name)
		if got != tt.want {
			t.Errorf("❌ IsExported(%q) = %v, want %v  ← Hint: check name[0] >= 'A' && name[0] <= 'Z'", tt.name, got, tt.want)
		} else {
			t.Logf("✅ IsExported(%q) = %v", tt.name, got)
		}
	}
}

func TestInitRan(t *testing.T) {
	log := GetInitLog()
	if len(log) == 0 {
		t.Error("❌ init() should have run and added to initLog")
	} else if log[0] != "packages_modules.init ran" {
		t.Errorf("❌ initLog[0] = %q, want 'packages_modules.init ran'", log[0])
	} else {
		t.Logf("✅ init() ran automatically: %q", log[0])
	}
}

func TestBlankImportPurpose(t *testing.T) {
	got := BlankImportPurpose()
	if got == "" {
		t.Error("❌ BlankImportPurpose() should return a non-empty explanation  ← Hint: it's about side effects / init()")
	} else {
		t.Logf("✅ BlankImportPurpose() = %q", got)
	}
}

func TestModulePath(t *testing.T) {
	got := ModulePath("github.com/user/myapp", "config")
	want := "github.com/user/myapp/config"
	if got != want {
		t.Errorf("❌ ModulePath = %q, want %q  ← Hint: return moduleName + \"/\" + subPackage", got, want)
	} else {
		t.Logf("✅ ModulePath(\"github.com/user/myapp\", \"config\") = %q", got)
	}
}
