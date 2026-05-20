package forgeconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsKnownLanguage(t *testing.T) {
	for _, lang := range KnownLanguages {
		if !IsKnownLanguage(lang) {
			t.Errorf("IsKnownLanguage(%q) should be true", lang)
		}
	}
	if IsKnownLanguage("unknown") {
		t.Error("IsKnownLanguage('unknown') should be false")
	}
	if IsKnownLanguage("") {
		t.Error("IsKnownLanguage('') should be false")
	}
}

func TestDetectLanguages(t *testing.T) {
	t.Run("go project detected from go.mod", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/test\ngo 1.21\n"), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 1 || langs[0] != "go" {
			t.Errorf("expected [go], got %v", langs)
		}
	})

	t.Run("javascript project detected from playwright dep", func(t *testing.T) {
		dir := t.TempDir()
		pkgJSON := `{"devDependencies":{"@playwright/test":"^1.0.0"}}`
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkgJSON), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		found := false
		for _, l := range langs {
			if l == "javascript" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected javascript in %v", langs)
		}
	})

	t.Run("rust project detected from Cargo.toml", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte("[package]\nname = \"test\"\n"), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 1 || langs[0] != "rust" {
			t.Errorf("expected [rust], got %v", langs)
		}
	})

	t.Run("java project detected from pom.xml", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project></project>"), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 1 || langs[0] != "java" {
			t.Errorf("expected [java], got %v", langs)
		}
	})

	t.Run("mobile project detected from android dir", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "android"), 0755); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 1 || langs[0] != "mobile" {
			t.Errorf("expected [mobile], got %v", langs)
		}
	})

	t.Run("python project detected from pytest in requirements.txt", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte("pytest>=7.0\n"), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 1 || langs[0] != "python" {
			t.Errorf("expected [python], got %v", langs)
		}
	})

	t.Run("python project detected from pytest in pyproject.toml", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("[tool.pytest]\n"), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 1 || langs[0] != "python" {
			t.Errorf("expected [python], got %v", langs)
		}
	})

	t.Run("empty project returns empty", func(t *testing.T) {
		dir := t.TempDir()
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 0 {
			t.Errorf("expected empty, got %v", langs)
		}
	})

	t.Run("deduplicates languages", func(t *testing.T) {
		dir := t.TempDir()
		// Both go.mod and android dir -> go + mobile, no duplicates
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(dir, "android"), 0755); err != nil {
			t.Fatal(err)
		}
		langs, err := DetectLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		seen := make(map[string]int)
		for _, l := range langs {
			seen[l]++
		}
		for l, c := range seen {
			if c > 1 {
				t.Errorf("language %q appears %d times", l, c)
			}
		}
	})
}

func TestReadLanguages(t *testing.T) {
	t.Run("config languages take priority", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}
		configYAML := "languages:\n  - go\n  - python\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configYAML), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := ReadLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 2 || langs[0] != "go" || langs[1] != "python" {
			t.Errorf("expected [go python], got %v", langs)
		}
	})

	t.Run("auto-detect when config has no languages", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644); err != nil {
			t.Fatal(err)
		}
		langs, err := ReadLanguages(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(langs) != 1 || langs[0] != "go" {
			t.Errorf("expected [go], got %v", langs)
		}
	})
}

func TestUnionLanguageInterfaces(t *testing.T) {
	t.Run("go gives api+cli", func(t *testing.T) {
		ifaces, err := UnionLanguageInterfaces([]string{"go"})
		if err != nil {
			t.Fatal(err)
		}
		if len(ifaces) != 2 || ifaces[0] != "api" || ifaces[1] != "cli" {
			t.Errorf("expected [api cli], got %v", ifaces)
		}
	})

	t.Run("javascript gives api+web-ui", func(t *testing.T) {
		ifaces, err := UnionLanguageInterfaces([]string{"javascript"})
		if err != nil {
			t.Fatal(err)
		}
		if len(ifaces) != 2 || ifaces[0] != "api" || ifaces[1] != "web-ui" {
			t.Errorf("expected [api web-ui], got %v", ifaces)
		}
	})

	t.Run("union deduplicates", func(t *testing.T) {
		ifaces, err := UnionLanguageInterfaces([]string{"go", "python"})
		if err != nil {
			t.Fatal(err)
		}
		seen := make(map[string]int)
		for _, i := range ifaces {
			seen[i]++
		}
		for k, v := range seen {
			if v > 1 {
				t.Errorf("interface %q duplicated", k)
			}
		}
	})

	t.Run("unknown language skipped", func(t *testing.T) {
		ifaces, err := UnionLanguageInterfaces([]string{"cobol"})
		if err != nil {
			t.Fatal(err)
		}
		if len(ifaces) != 0 {
			t.Errorf("expected empty for unknown, got %v", ifaces)
		}
	})

	t.Run("empty languages returns empty", func(t *testing.T) {
		ifaces, err := UnionLanguageInterfaces(nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(ifaces) != 0 {
			t.Errorf("expected empty, got %v", ifaces)
		}
	})
}

func TestReadInterfaces(t *testing.T) {
	t.Run("config interfaces take priority", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}
		configYAML := "interfaces:\n  - api\n  - cli\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configYAML), 0644); err != nil {
			t.Fatal(err)
		}
		ifaces, err := ReadInterfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ifaces) != 2 || ifaces[0] != "api" || ifaces[1] != "cli" {
			t.Errorf("expected [api cli], got %v", ifaces)
		}
	})

	t.Run("auto-detect from languages when no config interfaces", func(t *testing.T) {
		dir := t.TempDir()
		forgeDir := filepath.Join(dir, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}
		configYAML := "languages:\n  - go\n"
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configYAML), 0644); err != nil {
			t.Fatal(err)
		}
		ifaces, err := ReadInterfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ifaces) < 2 {
			t.Errorf("expected at least 2 interfaces for go, got %v", ifaces)
		}
	})
}
