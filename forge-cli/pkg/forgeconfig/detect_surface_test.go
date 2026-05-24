package forgeconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// --- AC1: Detects surface types from various manifest files ---

func TestDetectSurfaces_PackageJSON(t *testing.T) {
	tests := []struct {
		name     string
		deps     map[string]string // dependencies
		devDeps  map[string]string // devDependencies
		wantType string
	}{
		{
			name:     "react -> web",
			deps:     map[string]string{"react": "^18.0.0"},
			wantType: "web",
		},
		{
			name:     "vue -> web",
			deps:     map[string]string{"vue": "^3.0.0"},
			wantType: "web",
		},
		{
			name:     "svelte -> web",
			devDeps:  map[string]string{"svelte": "^4.0.0"},
			wantType: "web",
		},
		{
			name:     "express (no frontend) -> api",
			deps:     map[string]string{"express": "^4.18.0"},
			wantType: "api",
		},
		{
			name:     "fastify (no frontend) -> api",
			deps:     map[string]string{"fastify": "^4.0.0"},
			wantType: "api",
		},
		{
			name:     "commander (no frontend) -> cli",
			deps:     map[string]string{"commander": "^11.0.0"},
			wantType: "cli",
		},
		{
			name:     "yargs (no frontend) -> cli",
			deps:     map[string]string{"yargs": "^17.0.0"},
			wantType: "cli",
		},
		{
			name:     "blessed -> tui",
			deps:     map[string]string{"blessed": "^0.1.81"},
			wantType: "tui",
		},
		{
			name:     "ink -> tui",
			deps:     map[string]string{"ink": "^4.0.0"},
			wantType: "tui",
		},
		{
			name:     "react + express conflict -> web (higher priority)",
			deps:     map[string]string{"react": "^18.0.0", "express": "^4.18.0"},
			wantType: "web",
		},
		{
			name:     "react-native + react conflict -> mobile",
			deps:     map[string]string{"react-native": "^0.72.0", "react": "^18.0.0"},
			wantType: "mobile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writePackageJSON(t, dir, tt.deps, tt.devDeps)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 surface, got %d: %v", len(result), result)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

func TestDetectSurfaces_GoMod(t *testing.T) {
	tests := []struct {
		name     string
		require  []string // lines in require block
		wantType string
	}{
		{
			name:     "gin -> api",
			require:  []string{"github.com/gin-gonic/gin v1.9.0"},
			wantType: "api",
		},
		{
			name:     "echo -> api",
			require:  []string{"github.com/labstack/echo/v4 v4.11.0"},
			wantType: "api",
		},
		{
			name:     "cobra -> cli",
			require:  []string{"github.com/spf13/cobra v1.7.0"},
			wantType: "cli",
		},
		{
			name:     "bubbletea -> tui",
			require:  []string{"github.com/charmbracelet/bubbletea v0.24.0"},
			wantType: "tui",
		},
		{
			name:     "gin + cobra conflict -> api (higher priority)",
			require:  []string{"github.com/gin-gonic/gin v1.9.0", "github.com/spf13/cobra v1.7.0"},
			wantType: "api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeGoMod(t, dir, "example.com/test", tt.require)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 surface, got %d: %v", len(result), result)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

func TestDetectSurfaces_CargoToml(t *testing.T) {
	tests := []struct {
		name     string
		deps     []string // dependency lines
		wantType string
	}{
		{
			name:     "actix -> api",
			deps:     []string{"actix-web = \"4\""},
			wantType: "api",
		},
		{
			name:     "axum -> api",
			deps:     []string{"axum = \"0.7\""},
			wantType: "api",
		},
		{
			name:     "clap -> cli",
			deps:     []string{"clap = { version = \"4\", features = [\"derive\"] }"},
			wantType: "cli",
		},
		{
			name:     "ratatui -> tui",
			deps:     []string{"ratatui = \"0.25\""},
			wantType: "tui",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeCargoToml(t, dir, tt.deps)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 surface, got %d: %v", len(result), result)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

func TestDetectSurfaces_MobileProjects(t *testing.T) {
	t.Run("AndroidManifest.xml -> mobile", func(t *testing.T) {
		dir := t.TempDir()
		writeAndroidManifest(t, dir)

		result, err := DetectSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["."] != "mobile" {
			t.Errorf("expected 'mobile', got '%s'", result["."])
		}
	})

	t.Run("xcodeproj -> mobile", func(t *testing.T) {
		dir := t.TempDir()
		xcodeDir := filepath.Join(dir, "MyApp.xcodeproj")
		if err := os.MkdirAll(xcodeDir, 0755); err != nil {
			t.Fatal(err)
		}
		// Create project.pbxproj to make it look real
		if err := os.WriteFile(filepath.Join(xcodeDir, "project.pbxproj"), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}

		result, err := DetectSurfaces(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["."] != "mobile" {
			t.Errorf("expected 'mobile', got '%s'", result["."])
		}
	})
}

func TestDetectSurfaces_PyProject(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantType string
	}{
		{
			name:     "flask -> api",
			content:  "[project]\ndependencies = [\"flask\"]\n",
			wantType: "api",
		},
		{
			name:     "fastapi -> api",
			content:  "[project]\ndependencies = [\"fastapi\"]\n",
			wantType: "api",
		},
		{
			name:     "click -> cli",
			content:  "[project]\ndependencies = [\"click\"]\n",
			wantType: "cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result["."] != tt.wantType {
				t.Errorf("expected surfaces['.']='%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

// --- AC2: Workspace detection ---

func TestDetectSurfaces_WorkspacePnpm(t *testing.T) {
	dir := t.TempDir()

	// Root package.json (workspace config only)
	writePackageJSON(t, dir, map[string]string{"express": "^4.18.0"}, nil)
	// pnpm-workspace.yaml triggers workspace mode
	writeTestFile(t, dir, "pnpm-workspace.yaml", "packages:\n  - 'apps/*'\n  - 'packages/*'\n")

	// Subdir: apps/web has react
	webDir := filepath.Join(dir, "apps", "web")
	mkdirAll(t, webDir)
	writePackageJSON(t, webDir, map[string]string{"react": "^18.0.0"}, nil)

	// Subdir: apps/api has express
	apiDir := filepath.Join(dir, "apps", "api")
	mkdirAll(t, apiDir)
	writePackageJSON(t, apiDir, map[string]string{"express": "^4.18.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Root deps should be skipped; only subdirs detected
	if _, ok := result["."]; ok {
		t.Error("root deps should be skipped in workspace mode")
	}
	webKey := filepath.ToSlash(filepath.Join("apps", "web"))
	apiKey := filepath.ToSlash(filepath.Join("apps", "api"))
	if result[webKey] != "web" {
		t.Errorf("expected apps/web=web, got %v", result)
	}
	if result[apiKey] != "api" {
		t.Errorf("expected apps/api=api, got %v", result)
	}
}

func TestDetectSurfaces_WorkspacePackageJSON(t *testing.T) {
	dir := t.TempDir()

	// Root package.json with workspaces field
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*", "packages/*"})

	// Subdir
	cliDir := filepath.Join(dir, "packages", "cli")
	mkdirAll(t, cliDir)
	writePackageJSON(t, cliDir, map[string]string{"commander": "^11.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cliKey := filepath.ToSlash(filepath.Join("packages", "cli"))
	if result[cliKey] != "cli" {
		t.Errorf("expected packages/cli=cli, got %v", result)
	}
}

// --- AC3: Non-workspace root detected as "." ---

func TestDetectSurfaces_NonWorkspace_ScalarOutput(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Single type -> scalar form (key is ".")
	if len(result) != 1 {
		t.Fatalf("expected 1 surface, got %d", len(result))
	}
	if result["."] != "web" {
		t.Errorf("expected surfaces['.']='web', got '%s'", result["."])
	}
}

// --- AC4: Depth limit ---

func TestDetectSurfaces_DefaultDepth3(t *testing.T) {
	dir := t.TempDir()
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// Create at depth 3: apps/web/deep (3 levels from root)
	deepDir := filepath.Join(dir, "apps", "web", "deep")
	mkdirAll(t, deepDir)
	writePackageJSON(t, deepDir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deepKey := filepath.ToSlash(filepath.Join("apps", "web", "deep"))
	if result[deepKey] != "web" {
		t.Errorf("expected to detect at depth 3, got %v", result)
	}
}

func TestDetectSurfaces_DepthExceeded(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "2")

	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// Create at depth 3: should NOT be detected
	deepDir := filepath.Join(dir, "apps", "web", "deep")
	mkdirAll(t, deepDir)
	writePackageJSON(t, deepDir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deepKey := filepath.ToSlash(filepath.Join("apps", "web", "deep"))
	if _, ok := result[deepKey]; ok {
		t.Errorf("should NOT detect at depth 3 with FORGE_DETECT_DEPTH=2, got %v", result)
	}
}

func TestDetectSurfaces_DepthZeroError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "0")

	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	_, err := DetectSurfaces(dir)
	if err == nil {
		t.Fatal("expected error for FORGE_DETECT_DEPTH=0")
	}
	if _, ok := err.(*ErrInvalidDepth); !ok {
		t.Errorf("expected ErrInvalidDepth, got %T: %v", err, err)
	}
}

func TestDetectSurfaces_DepthNegativeError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "-1")

	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	_, err := DetectSurfaces(dir)
	if err == nil {
		t.Fatal("expected error for FORGE_DETECT_DEPTH=-1")
	}
}

func TestDetectSurfaces_DepthEnvOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("FORGE_DETECT_DEPTH", "1")

	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// apps/web is depth 2 from root, exceeds FORGE_DETECT_DEPTH=1
	skipDir := filepath.Join(dir, "apps", "web")
	mkdirAll(t, skipDir)
	writePackageJSON(t, skipDir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	skipKey := filepath.ToSlash(filepath.Join("apps", "web"))
	if _, ok := result[skipKey]; ok {
		t.Errorf("should NOT detect at depth 2 with FORGE_DETECT_DEPTH=1, got %v", result)
	}
}

// --- AC5: Exclusion dirs ---

func TestDetectSurfaces_ExclusionDirs(t *testing.T) {
	dir := t.TempDir()
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	// node_modules should be skipped
	nmDir := filepath.Join(dir, "node_modules", "pkg")
	mkdirAll(t, nmDir)
	writePackageJSON(t, nmDir, map[string]string{"react": "^18.0.0"}, nil)

	// .git should be skipped
	gitDir := filepath.Join(dir, ".git", "subdir")
	mkdirAll(t, gitDir)
	writePackageJSON(t, gitDir, map[string]string{"express": "^4.18.0"}, nil)

	// dist should be skipped
	distDir := filepath.Join(dir, "dist", "pkg")
	mkdirAll(t, distDir)
	writePackageJSON(t, distDir, map[string]string{"commander": "^11.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) > 0 {
		t.Errorf("expected 0 surfaces (all dirs excluded), got %d: %v", len(result), result)
	}
}

func TestDetectSurfaces_ExclusionDirList(t *testing.T) {
	expectedExclusions := []string{
		"node_modules", ".git", "vendor", "dist",
		"build", "__pycache__", ".next", "target",
	}

	for _, d := range expectedExclusions {
		if !isExcludedDir(d) {
			t.Errorf("expected %q to be excluded", d)
		}
	}
}

// --- AC6: Signal conflict resolution ---

func TestDetectSurfaces_ConflictResolution(t *testing.T) {
	tests := []struct {
		name     string
		deps     map[string]string
		wantType string
	}{
		{
			name:     "web > api (react + express)",
			deps:     map[string]string{"react": "^18.0.0", "express": "^4.18.0"},
			wantType: "web",
		},
		{
			name:     "web > cli (react + commander)",
			deps:     map[string]string{"react": "^18.0.0", "commander": "^11.0.0"},
			wantType: "web",
		},
		{
			name:     "api > cli (express + commander)",
			deps:     map[string]string{"express": "^4.18.0", "commander": "^11.0.0"},
			wantType: "api",
		},
		{
			name:     "api > tui (express + blessed)",
			deps:     map[string]string{"express": "^4.18.0", "blessed": "^0.1.81"},
			wantType: "api",
		},
		{
			name:     "cli > tui (commander + blessed)",
			deps:     map[string]string{"commander": "^11.0.0", "blessed": "^0.1.81"},
			wantType: "cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writePackageJSON(t, dir, tt.deps, nil)

			result, err := DetectSurfaces(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result["."] != tt.wantType {
				t.Errorf("expected '%s', got '%s'", tt.wantType, result["."])
			}
		})
	}
}

// --- AC7: Single-type -> scalar; multi-type -> map ---

func TestDetectSurfaces_SingleTypeScalar(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"react": "^18.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 surface, got %d", len(result))
	}
	if result["."] != "web" {
		t.Errorf("expected scalar form with key '.', got %v", result)
	}
}

func TestDetectSurfaces_MultiTypeMap(t *testing.T) {
	dir := t.TempDir()
	writePackageJSONWithWorkspaces(t, dir, []string{"apps/*"})

	webDir := filepath.Join(dir, "apps", "frontend")
	mkdirAll(t, webDir)
	writePackageJSON(t, webDir, map[string]string{"react": "^18.0.0"}, nil)

	apiDir := filepath.Join(dir, "apps", "backend")
	mkdirAll(t, apiDir)
	writePackageJSON(t, apiDir, map[string]string{"express": "^4.18.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 surfaces, got %d: %v", len(result), result)
	}

	frontendKey := filepath.ToSlash(filepath.Join("apps", "frontend"))
	backendKey := filepath.ToSlash(filepath.Join("apps", "backend"))
	if result[frontendKey] != "web" {
		t.Errorf("expected apps/frontend=web, got %v", result)
	}
	if result[backendKey] != "api" {
		t.Errorf("expected apps/backend=api, got %v", result)
	}
}

// --- Edge cases ---

func TestDetectSurfaces_NoManifests(t *testing.T) {
	dir := t.TempDir()

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for empty dir, got %v", result)
	}
}

func TestDetectSurfaces_UnknownDepsIgnored(t *testing.T) {
	dir := t.TempDir()
	writePackageJSON(t, dir, map[string]string{"lodash": "^4.17.0", "axios": "^1.0.0"}, nil)

	result, err := DetectSurfaces(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result for unknown deps, got %v", result)
	}
}

// --- Helper functions ---

func writePackageJSON(t *testing.T, dir string, deps, devDeps map[string]string) {
	t.Helper()
	type pkg struct {
		Dependencies    map[string]string `json:"dependencies,omitempty"`
		DevDependencies map[string]string `json:"devDependencies,omitempty"`
	}
	p := pkg{
		Dependencies:    deps,
		DevDependencies: devDeps,
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func writePackageJSONWithWorkspaces(t *testing.T, dir string, workspaces []string) {
	t.Helper()
	type pkg struct {
		Workspaces []string `json:"workspaces"`
	}
	p := pkg{Workspaces: workspaces}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeGoMod(t *testing.T, dir, module string, require []string) {
	t.Helper()
	content := "module " + module + "\n\ngo 1.21\n"
	if len(require) > 0 {
		content += "\nrequire (\n"
		for _, r := range require {
			content += "\t" + r + "\n"
		}
		content += ")\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeCargoToml(t *testing.T, dir string, deps []string) {
	t.Helper()
	content := "[package]\nname = \"test\"\nversion = \"0.1.0\"\n\n[dependencies]\n"
	for _, d := range deps {
		content += d + "\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeAndroidManifest(t *testing.T, dir string) {
	t.Helper()
	androidDir := filepath.Join(dir, "app", "src", "main")
	mkdirAll(t, androidDir)
	content := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<manifest xmlns:android=\"http://schemas.android.com/apk/res/android\"\n    package=\"com.example.app\">\n</manifest>"
	if err := os.WriteFile(filepath.Join(androidDir, "AndroidManifest.xml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}
