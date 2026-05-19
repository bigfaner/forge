package profile

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestGetBuiltinFramework(t *testing.T) {
	t.Run("go-testing", func(t *testing.T) {
		fw, err := GetBuiltinFramework("go-testing")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "go-testing" {
			t.Errorf("Name = %q, want %q", fw.Name, "go-testing")
		}
		if fw.TestFunctionPattern != "func Test*" {
			t.Errorf("TestFunctionPattern = %q, want %q", fw.TestFunctionPattern, "func Test*")
		}
		if fw.FilePattern != "*_test.go" {
			t.Errorf("FilePattern = %q, want %q", fw.FilePattern, "*_test.go")
		}
		if fw.LanguageHint != "go" {
			t.Errorf("LanguageHint = %q, want %q", fw.LanguageHint, "go")
		}
	})

	t.Run("pytest", func(t *testing.T) {
		fw, err := GetBuiltinFramework("pytest")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.TestFunctionPattern != "def test_*" {
			t.Errorf("TestFunctionPattern = %q, want %q", fw.TestFunctionPattern, "def test_*")
		}
		if fw.FilePattern != "test_*.py" {
			t.Errorf("FilePattern = %q, want %q", fw.FilePattern, "test_*.py")
		}
	})

	t.Run("mocha", func(t *testing.T) {
		fw, err := GetBuiltinFramework("mocha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.TestFunctionPattern != "describe/it" {
			t.Errorf("TestFunctionPattern = %q, want %q", fw.TestFunctionPattern, "describe/it")
		}
		if fw.FilePattern != "*.spec.ts" {
			t.Errorf("FilePattern = %q, want %q", fw.FilePattern, "*.spec.ts")
		}
	})

	t.Run("unknown returns error", func(t *testing.T) {
		_, err := GetBuiltinFramework("unknown-framework")
		if err == nil {
			t.Fatal("expected error for unknown framework")
		}
	})
}

func TestIsKnownFramework(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"go-testing", true},
		{"pytest", true},
		{"mocha", true},
		{"junit5", true},
		{"rust-test", true},
		{"maestro", true},
		{"unknown", false},
		{"", false},
		{"Go-Testing", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKnownFramework(tt.name); got != tt.expected {
				t.Errorf("IsKnownFramework(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestKnownFrameworkNames(t *testing.T) {
	names := KnownFrameworkNames()
	expected := []string{"go-testing", "junit5", "maestro", "mocha", "pytest", "rust-test"}
	if len(names) != len(expected) {
		t.Fatalf("KnownFrameworkNames() = %d, want %d", len(names), len(expected))
	}
	for _, name := range expected {
		if !slices.Contains(names, name) {
			t.Errorf("expected framework %q not found", name)
		}
	}
	// Verify sorted
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("KnownFrameworkNames() not sorted: %q before %q", names[i-1], names[i])
		}
	}
}

func TestDefaultFrameworkForLanguage(t *testing.T) {
	tests := []struct {
		lang       string
		wantFw     string
		wantExists bool
	}{
		{"go", "go-testing", true},
		{"python", "pytest", true},
		{"javascript", "mocha", true},
		{"java", "junit5", true},
		{"rust", "rust-test", true},
		{"mobile", "maestro", true},
		{"unknown", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			fw, ok := DefaultFrameworkForLanguage(tt.lang)
			if ok != tt.wantExists {
				t.Errorf("DefaultFrameworkForLanguage(%q) exists = %v, want %v", tt.lang, ok, tt.wantExists)
			}
			if fw != tt.wantFw {
				t.Errorf("DefaultFrameworkForLanguage(%q) = %q, want %q", tt.lang, fw, tt.wantFw)
			}
		})
	}
}

func TestResolveTestFramework(t *testing.T) {
	setupConfig := func(t *testing.T, content string) string {
		t.Helper()
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		return dir
	}

	t.Run("explicit test-framework in config", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\ntest-framework: pytest\n")

		fw, err := ResolveTestFramework(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "pytest" {
			t.Errorf("expected pytest, got %q", fw.Name)
		}
		if fw.TestFunctionPattern != "def test_*" {
			t.Errorf("expected 'def test_*', got %q", fw.TestFunctionPattern)
		}
	})

	t.Run("no test-framework defaults to language framework", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - go\n")

		fw, err := ResolveTestFramework(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "go-testing" {
			t.Errorf("expected go-testing, got %q", fw.Name)
		}
	})

	t.Run("no test-framework javascript defaults to mocha", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - javascript\n")

		fw, err := ResolveTestFramework(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "mocha" {
			t.Errorf("expected mocha, got %q", fw.Name)
		}
	})

	t.Run("no config returns empty framework", func(t *testing.T) {
		dir := t.TempDir()

		fw, err := ResolveTestFramework(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "" {
			t.Errorf("expected empty framework, got %q", fw.Name)
		}
	})

	t.Run("config with no languages returns empty framework", func(t *testing.T) {
		dir := setupConfig(t, "project-type: backend\n")

		fw, err := ResolveTestFramework(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "" {
			t.Errorf("expected empty framework, got %q", fw.Name)
		}
	})

	t.Run("custom framework name not in builtins", func(t *testing.T) {
		dir := setupConfig(t, "test-framework: my-custom-fw\n")

		fw, err := ResolveTestFramework(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "my-custom-fw" {
			t.Errorf("expected my-custom-fw, got %q", fw.Name)
		}
		// Custom framework has no pattern info — only the name from config
		if fw.TestFunctionPattern != "" {
			t.Errorf("expected empty TestFunctionPattern for custom framework, got %q", fw.TestFunctionPattern)
		}
	})

	t.Run("first language used for default when multiple languages", func(t *testing.T) {
		dir := setupConfig(t, "languages:\n  - python\n  - go\n")

		fw, err := ResolveTestFramework(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw.Name != "pytest" {
			t.Errorf("expected pytest (first language python), got %q", fw.Name)
		}
	})
}

func TestRegisterCustomFramework(t *testing.T) {
	custom := FrameworkInfo{
		Name:                "custom-fw",
		TestFunctionPattern: "custom_test_*",
		FilePattern:         "custom_*.ext",
		LanguageHint:        "custom-lang",
	}

	// Register
	RegisterCustomFramework(custom)

	// Verify it can be retrieved
	fw, err := GetBuiltinFramework("custom-fw")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fw.Name != "custom-fw" {
		t.Errorf("Name = %q, want %q", fw.Name, "custom-fw")
	}
	if fw.TestFunctionPattern != "custom_test_*" {
		t.Errorf("TestFunctionPattern = %q, want %q", fw.TestFunctionPattern, "custom_test_*")
	}
	if fw.FilePattern != "custom_*.ext" {
		t.Errorf("FilePattern = %q, want %q", fw.FilePattern, "custom_*.ext")
	}
	if fw.LanguageHint != "custom-lang" {
		t.Errorf("LanguageHint = %q, want %q", fw.LanguageHint, "custom-lang")
	}

	// Verify IsKnownFramework recognizes it
	if !IsKnownFramework("custom-fw") {
		t.Error("IsKnownFramework should recognize custom registered framework")
	}

	// Clean up — remove the custom framework to not affect other tests
	delete(builtinFrameworks, "custom-fw")
}

func TestRegisterCustomFrameworkOverrides(t *testing.T) {
	// Register a custom framework that overrides a built-in name
	custom := FrameworkInfo{
		Name:                "go-testing",
		TestFunctionPattern: "custom_go_pattern",
		FilePattern:         "custom_*.go",
	}

	RegisterCustomFramework(custom)

	fw, err := GetBuiltinFramework("go-testing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fw.TestFunctionPattern != "custom_go_pattern" {
		t.Errorf("expected custom override, got %q", fw.TestFunctionPattern)
	}

	// Restore original
	builtinFrameworks["go-testing"] = FrameworkInfo{
		Name:                "go-testing",
		TestFunctionPattern: "func Test*",
		FilePattern:         "*_test.go",
		LanguageHint:        "go",
	}
}
